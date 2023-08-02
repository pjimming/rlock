package rlock

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pjimming/rlock/constants"
	"github.com/pjimming/rlock/logx"
	"github.com/pjimming/rlock/utils"

	"github.com/go-redis/redis/v8"
)

var (
	once sync.Once
	rc   *redis.Client
)

type RLock struct {
	key    string        // lock key
	token  string        // lock token, to authority validate
	client *redis.Client // redis client pool
	lockOptions

	runningDog int32              // watchdog running status
	stopDog    context.CancelFunc // stop watchdog

	ctx    context.Context
	logger *logx.Logger // log
}

// NewRLock new a redis lock.
// You can set params with set function.
func NewRLock(op RedisClientOptions, key string) (rLock *RLock) {
	if key == "" {
		key = utils.GenerateRandomString(10)
	}

	once.Do(func() {
		rc = redis.NewClient(&redis.Options{
			Addr:     op.Addr,
			Password: op.Password,

			//连接池容量及闲置连接数量
			PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
			MinIdleConns: 10, // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量

			// 超时
			DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒。
			ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
			WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时
			PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

			// 闲置连接检查包括IdleTimeout，MaxConnAge
			IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
			IdleTimeout:        5 * time.Minute,  // 闲置超时，默认5分钟，-1表示取消闲置超时检查
			MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

			// 命令执行失败时的重试策略
			MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
			MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
		})
	})

	if err := rc.Ping(rc.Context()).Err(); err != nil {
		log.Printf("redis client ping fail, error: %v", err)
		return nil
	}

	rLock = &RLock{
		key:    key,
		token:  utils.GenerateToken(),
		client: rc,
		lockOptions: lockOptions{
			blockWaitingSecond: 60,
			expireSeconds:      30,
			watchdogSwitch:     false,
		},
		runningDog: 0,
		stopDog:    nil,
		ctx:        context.Background(),
		logger:     logx.NewLogger(),
	}

	return
}

// Lock try to acquire lock, until acquire lock or blocking timeout. Returns:
//
// 1. ttl < 0 error;
//
// 2. ttl = 0 success;
//
// 3. ttl > 0 acquire by others
func (l *RLock) Lock() int64 {
	ttl, err := l.tryLock()
	if err != nil {
		return -1
	}

	if vv, ok := ttl.(int64); ok && vv == int64(0) {
		return 0
	}

	// span
	ttl, err = l.span()
	if err != nil {
		return -1
	}
	if vv, ok := ttl.(int64); ok {
		return vv
	}
	l.logger.Debug("Lock ttl interface trans int64 fail")
	return -1
}

// TryLock try to acquire lock in once, support reentrant. Returns:
//
// 1. ttl < 0 error;
//
// 2. ttl = 0 success;
//
// 3. ttl > 0 acquire by others.
func (l *RLock) TryLock() int64 {
	ttl, err := l.tryLock()
	if err != nil {
		return -1
	}
	if vv, ok := ttl.(int64); ok {
		return vv
	}
	l.logger.Debug("TryLock ttl interface trans int64 fail")
	return -1
}

// UnLock try to release lock, support reentrant. Returns:
//
// 1. res < 0 error or release other's lock;
//
// 2. res = 0 release lock success, but still hold lock because of reentry;
//
// 3. res = 1 release lock success absolutely.
func (l *RLock) UnLock() int64 {
	res, err := l.releaseLock()
	if err != nil {
		return -1
	}
	if vv, ok := res.(int64); ok {
		return vv
	}
	l.logger.Debug("UnLock ttl interface trans int64 fail")
	return -1
}

// tryLock try to acquire lock.
// if ttl == 0 means acquire lock successfully.
func (l *RLock) tryLock() (ttl interface{}, err error) {
	defer func() {
		if err != nil || l.runningDog == int32(1) {
			return
		}

		if vv, ok := ttl.(int64); ok && vv == int64(0) {
			l.watchdog()
		}
	}()

	if ttl, err = l.client.
		Eval(l.ctx, constants.LockLua, []string{l.key}, l.token, l.expireSeconds).
		Result(); err != nil {
		l.logger.Errorf("try lock fail, error: %v", err)
		return -1, err
	}
	return
}

// span try to acquire lock. If acquire unsuccessful,
// will blocking poll to acquire lock until context timeout or blocking timeout.
// Returns: ttl and error. If ttl == 0 means acquire lock successfully.
func (l *RLock) span() (ttl interface{}, err error) {
	timeoutCh := time.After(time.Duration(l.blockWaitingSecond) * time.Second)
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-l.ctx.Done():
			return -1, fmt.Errorf("span fail, ctx timeout, error: %v", l.ctx.Err())
		case <-timeoutCh:
			return l.tryLock()
		default:
		}

		ttl, err = l.tryLock()
		if err != nil {
			return -1, err
		}

		if vv, ok := ttl.(int64); ok && vv == int64(0) {
			return ttl, nil
		}
	}

	// will never reach
	return -1, nil
}

// releaseLock try to release lock.
// If res == 0, it means release lock successfully, but still hold lock.
// If res == 1, it means release lock absolutely.
// If res == -1, it means release other's lock or error.
func (l *RLock) releaseLock() (res interface{}, err error) {
	defer func() {
		if err != nil || res != int64(1) {
			return
		}

		// release lock absolutely
		if l.stopDog != nil {
			l.stopDog()
		}
	}()

	if res, err = l.client.
		Eval(l.ctx, constants.UnLockLua, []string{l.key}, l.token, l.expireSeconds).
		Result(); err != nil {
		l.logger.Errorf("release lock fail, error: %v", err)
		return -1, err
	}
	return
}

// delayExpire try to delay lock expire time.
func (l *RLock) delayExpire() (res interface{}, err error) {
	if res, err = l.client.
		Eval(l.ctx, constants.DelayExpireLua, []string{l.key}, l.token, l.expireSeconds).
		Result(); err != nil {
		l.logger.Errorf("delay expire fail, error: %v", err)
		return -1, err
	}
	return
}

func (l *RLock) Key() string {
	return l.key
}

func (l *RLock) SetKey(key string) *RLock {
	if key == "" {
		return l
	}
	l.key = key
	return l
}

func (l *RLock) Token() string {
	return l.token
}

func (l *RLock) SetToken(token string) *RLock {
	if token == "" {
		return l
	}
	l.token = token
	return l
}

func (l *RLock) BlockWaitingSecond() int64 {
	return l.blockWaitingSecond
}

func (l *RLock) SetBlockWaitingSecond(blockWaitingSecond int64) *RLock {
	l.blockWaitingSecond = blockWaitingSecond
	return l
}

func (l *RLock) ExpireSeconds() int64 {
	return l.expireSeconds
}

func (l *RLock) SetExpireSeconds(expireSeconds int64) *RLock {
	l.expireSeconds = expireSeconds
	return l
}

func (l *RLock) WatchdogSwitch() bool {
	return l.watchdogSwitch
}

func (l *RLock) SetWatchdogSwitch(watchdogSwitch bool) *RLock {
	l.watchdogSwitch = watchdogSwitch
	return l
}
