package RedisLock

import (
	"errors"
)

func errInvalidRedisLockType() error {
	return errors.New(ErrInvalidRLockType)
}
