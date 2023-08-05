package rlock

import "time"

type RedisClientOptions struct {
	Addr     string // redis address
	Password string // redis password
}

type LockOptions struct {
	blockWaitingTime time.Duration // blocking timeout time, default 60s.
	expireTime       time.Duration // key expire time, default 30s.
	watchdogSwitch   bool          // watchdog on/off, default false.
}

type LockOption func(*LockOptions)

func WithBlockWaitingTime(blockWaitingTime time.Duration) LockOption {
	return func(o *LockOptions) {
		o.blockWaitingTime = blockWaitingTime
	}
}

func WithExpireTime(expireTime time.Duration) LockOption {
	return func(o *LockOptions) {
		o.expireTime = expireTime
	}
}

func WithWatchDogSwitch(watchdogSwitch bool) LockOption {
	return func(o *LockOptions) {
		o.watchdogSwitch = watchdogSwitch
	}
}
