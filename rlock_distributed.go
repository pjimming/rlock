package RedisLock

import (
	"time"
)

type RLockDistributed struct {
	client RedisClient
}

func newRLockDistributed() (*RLockDistributed, error) {
	return &RLockDistributed{client: rc}, nil
}

func (l *RLockDistributed) TryLock(key, value string, expiration time.Duration) (bool, error) {
	return l.client.SetNX(l.client.Context(), key, value, expiration).Result()
}

func (l *RLockDistributed) ReleaseLock(key, value string) (bool, error) {
	result, err := l.client.Eval(l.client.Context(), l.releaseLockLua(), []string{key}, value).Result()
	return result == int64(1), err
}

func (l *RLockDistributed) HoldLock(key, value string) bool {
	return l.client.Get(l.client.Context(), key).Val() == value
}

func (l *RLockDistributed) releaseLockLua() string {
	return `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else 
	return 0
end
`
}
