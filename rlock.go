package rlock

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type RLocker interface {
	TryLock(key, value string, expiration time.Duration) (bool, error)
	ReleaseLock(key, value string) (bool, error)
	HoldLock(key, value string) bool
}

func NewRedisLock(param Param) (rLock RLocker, err error) {
	once.Do(func() {
		if len(param.Addr) > 1 {
			rc = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:       param.Addr,
				Password:    param.Password,
				DialTimeout: param.Timeout * time.Millisecond,
			})
		} else {
			rc = redis.NewClient(&redis.Options{
				Addr:        param.Addr[0],
				Password:    param.Password,
				DialTimeout: param.Timeout * time.Millisecond,
			})
		}
	})

	if err = rc.Ping(rc.Context()).Err(); err != nil {
		dlog.Errorf("redis client ping fail, %v", err)
		return nil, err
	}

	switch param.Type {
	case TypeDistributed, "":
		rLock, err = newRLockDistributed()
	case TypeReentry:
		rLock, err = newRLockReentry()
	case TypeSpin:
	case TypeAsynRenew:
	default:
		err = errInvalidRedisLockType()
	}

	return
}
