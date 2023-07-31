package rlock

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// common variable
var (
	once sync.Once
	rc   *redis.Client
	dlog *DLog
)

// generateToken 生成token
func generateToken() string {
	return fmt.Sprintf("%s_%s", getCurrentProcessID(), getCurrentGoroutineID())
}

// getCurrentProcessID 获取当前进程ID
func getCurrentProcessID() string {
	return strconv.Itoa(os.Getpid())
}

// getCurrentGoroutineID 获取当前的协程ID
func getCurrentGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
}

// RLock setter and getter

func (l *RLock) Key() string {
	return l.key
}

func (l *RLock) SetKey(key string) {
	l.key = key
}

func (l *RLock) Token() string {
	return l.token
}

func (l *RLock) SetToken(token string) {
	l.token = token
}

func (l *RLock) IsReentry() bool {
	return l.isReentry
}

func (l *RLock) SetIsReentry(isReentry bool) {
	l.isReentry = isReentry
}

func (l *RLock) IsBlock() bool {
	return l.isBlock
}

func (l *RLock) SetIsBlock(isBlock bool) {
	l.isBlock = isBlock
}

func (l *RLock) BlockWaitingSecond() int64 {
	return l.blockWaitingSecond
}

func (l *RLock) SetBlockWaitingSecond(blockWaitingSecond int64) {
	l.blockWaitingSecond = blockWaitingSecond
}

func (l *RLock) ExpireSeconds() int64 {
	return l.expireSeconds
}

func (l *RLock) SetExpireSeconds(expireSeconds int64) {
	l.expireSeconds = expireSeconds
}

func (l *RLock) WatchdogSwitch() bool {
	return l.watchdogSwitch
}

func (l *RLock) SetWatchdogSwitch(watchdogSwitch bool) {
	l.watchdogSwitch = watchdogSwitch
}
