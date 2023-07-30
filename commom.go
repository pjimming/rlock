package rlock

import (
	"sync"
)

// common variable
var (
	once sync.Once
	rc   RedisClient
	dlog *DLog
)
