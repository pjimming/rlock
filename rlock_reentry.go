package rlock

import (
	"time"
)

type RLockReentry struct {
	client RedisClient
}

func newRLockReentry() (*RLockReentry, error) {
	return &RLockReentry{client: rc}, nil
}

func (l *RLockReentry) TryLock(key, value string, expiration time.Duration) (bool, error) {
	result, err := l.client.Eval(l.client.Context(), l.tryLockLua(), []string{key}, value, expiration).Result()
	return result == int64(1), err
}

func (l *RLockReentry) ReleaseLock(key, value string) (bool, error) {
	result, err := l.client.Eval(l.client.Context(), l.releaseLockLua(), []string{key}, value).Result()
	return result == int64(1), err
}

func (l *RLockReentry) HoldLock(key, value string) bool {
	result, err := l.client.Eval(l.client.Context(), l.holdLockLua(), []string{key}, value).Result()
	if err != nil {
		dlog.Errorf("check hold lock, %v", err)
		return false
	}
	return result == int64(1)
}

func (l *RLockReentry) tryLockLua() string {
	return `
local lockKey = KEYS[1]
local lockToken = ARGV[1]
local lockTimeout = tonumber(ARGV[2])

local currentToken = redis.call("HGET", lockKey, "token")

if lockToken == currentToken then
	redis.call("HINCRBY", lockKey, "count", 1)
	redis.call("EXPIRE", lockKey, lockTimeout)
	return 1
elseif not currentToken then
	redis.call("HMSET", lockKey, "token", lockToken, "count", 1)
	redis.call("EXPIRE", lockKey, lockTimeout)
	return 1
else
	return 0
end
`
}

func (l *RLockReentry) releaseLockLua() string {
	return `
local lockKey = KEYS[1]
local lockToken = ARGV[1]

local currentToken = redis.call("HGET", lockKey, "token")

if lockToken == currentToken then
	local counter = tonumber(redis.call("HGET", lockKey, "count") or 0)
	if counter > 0 then
		counter = counter - 1
		if counter > 0 then
			redis.call("HSET", lockKey, "count", counter)
			return 0
		else
			redis.call("DEL", lockKey)
			return 1
		end
	else
		redis.call("DEL", lockKey)
		return 1
	end
else
	return 0
end
`
}

func (l *RLockReentry) holdLockLua() string {
	return `
local lockKey = KEYS[1]
local lockToken = ARGV[1]

local currentToken = redis.call("HGET", lockKey, "token")

if currentToken == lockToken then
	return 1
else
	return 0
end
`
}
