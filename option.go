package rlock

type SingleClientOptions struct {
	// redis distributedParam
	Addr     string
	Password string
}

type lockOptions struct {
	isReentry          bool
	isBlock            bool
	blockWaitingSecond int64
	expireSeconds      int64
	watchdogSwitch     bool
}
