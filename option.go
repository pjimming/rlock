package rlock

type RedisClientOptions struct {
	// redis address
	Addr string
	// redis password
	Password string
}

type lockOptions struct {
	isReentry          bool
	blockWaitingSecond int64
	expireSeconds      int64
	watchdogSwitch     bool
}