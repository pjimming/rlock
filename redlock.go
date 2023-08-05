package rlock

import (
	"errors"
	"time"
)

type RedLock struct {
	locks                 []*RLock
	maxSingleNodeWaitTime time.Duration // max try lock wait time.
}

// NewRedLock new a RedLock from multi redis servers.
//
// It is required that the cumulative timeout threshold of all nodes is
// less than one-tenth of the distributed lock expiration time.
func NewRedLock(rcs []RedisClientOptions, key string, expireTime time.Duration) (redLock *RedLock, err error) {
	if key == "" {
		key = generateRandomString(10)
	}

	for _, op := range rcs {
		rlock := NewRLock(op, key)

		if rlock != nil {
			redLock.locks = append(redLock.locks, rlock.
				SetToken(key+"_token").
				SetWatchdogSwitch(true).
				SetExpireTime(expireTime))
		}
	}

	if len(redLock.locks) < 3 {
		return nil, errors.New("new redlock fail, locks count less than 3")
	}

	redLock.maxSingleNodeWaitTime = expireTime / time.Duration(10*len(redLock.locks))

	return
}

// TryLock try to acquire lock.
//
// If RedLock gets lock count greater than half of locks,
// it means acquire lock successfully.
func (l *RedLock) TryLock() bool {
	successCnt := 0
	for _, lock := range l.locks {
		start := time.Now()
		ttl := lock.TryLock()
		cost := time.Since(start)
		if ttl == int64(0) && cost <= l.maxSingleNodeWaitTime {
			successCnt++
		}
	}

	return successCnt >= (len(l.locks)>>1 + 1)
}

// UnLock release lock.
func (l *RedLock) UnLock() {
	for _, lock := range l.locks {
		lock.UnLock()
	}
}
