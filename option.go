package rlock

import "time"

type RedisClientOptions struct {
	Addr     string // redis address
	Password string // redis password
}

type lockOptions struct {
	blockWaitingTime time.Duration // blocking timeout time, default 60s.
	expireTime       time.Duration // key expire time, default 30s.
	watchdogSwitch   bool          // watchdog on/off, default false.
}

type redLockOptions struct {
	maxSingleNodeWaitTime time.Duration // max try lock wait time.
}
