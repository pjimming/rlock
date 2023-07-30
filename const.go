package rlock

// RedisLock Type
const (
	TypeDistributed = "distributed"
	TypeReentry     = "reentry"
	TypeSpin        = "spin"
	TypeAsynRenew   = "asyn_renew"
)

// Error Message
const (
	ErrInvalidRLockType = "invalid redis lock type"
)
