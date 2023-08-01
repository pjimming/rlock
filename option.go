package rlock

type RedisClientOptions struct {
	Addr     string // redis address
	Password string // redis password
}

type lockOptions struct {
	blockWaitingSecond int64 // blocking timeout time
	expireSeconds      int64 // key expire time
	watchdogSwitch     bool  // watchdog on/off
}
