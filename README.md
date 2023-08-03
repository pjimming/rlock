# RLock

A distributed lock base on Redis achieved by Golang.

- [简体中文](./README_ZH.md)

---

## Status
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Usage
```shell
go get -u github.com/pjimming/rlock
```

## Achieve
- Mutual exclusion: Redis distributed lock can ensure that only one client can acquire the lock at the same time, realizing mutual exclusion between threads.
- Security: Redis distributed locks use atomic operations, which can ensure the security of locks under concurrent conditions and avoid problems such as data competition and deadlocks.
- Lock timeout: In order to avoid deadlock caused by a failure of a certain client after acquiring the lock, the Redis distributed lock can set the lock timeout period, and the lock will be released automatically when the timeout is exceeded.
- Reentrancy: Redis distributed locks can support the same client to acquire the same lock multiple times, avoiding deadlocks in nested calls.
- High performance: Redis is an in-memory database with high read and write performance, enabling fast locking and unlocking operations.
- Atomicity: The locking and unlocking operations of Redis distributed locks use atomic commands, which can ensure the atomicity of operations and avoid competition problems under concurrency.

## Quick Start
```go
package test

import (
	"testing"
	"time"

	"github.com/pjimming/rlock"
	"github.com/pjimming/rlock/utils"

	"github.com/stretchr/testify/assert"
)

var op = rlock.RedisClientOptions{
	Addr:     "127.0.0.1:6379",
	Password: "",
}

func TestLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)

	ttl := l.Lock()
	ast.Equal(int64(0), ttl)

	l.UnLock()
}

func TestTryLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)

	ttl := l.TryLock()
	t.Log("ttl:", ttl)
	ast.Equal(int64(0), ttl)

	l.UnLock()
}

func TestDelayExpire(t *testing.T) {
	ast := assert.New(t)
	key := utils.GenerateRandomString(8)

	l1 := rlock.NewRLock(op, key).
		SetToken(key + "111").
		SetExpireSeconds(5).
		SetWatchdogSwitch(true)

	l2 := rlock.NewRLock(op, key).
		SetToken(key + "222").
		SetBlockWaitingSecond(20)

	t.Log("l1:", l1.Key(), l1.Token())
	t.Log("l2:", l2.Key(), l2.Token())

	l1.Lock()

	start := time.Now()
	l2.TryLock()
	t.Log("l2 TryLock cost:", time.Now().Sub(start).String())

	t.Log("l1 start sleep...", time.Now().Sub(start).String())
	time.Sleep(time.Second * 10)
	ast.Equal(int64(1), l1.UnLock())
	t.Log("l1 unlock", time.Now().Sub(start).String())

	l2.Lock()
	t.Log("l2 Lock cost:", time.Now().Sub(start).String())

	ast.Equal(int64(1), l2.UnLock())
}
```

## Lua Scripts
> Hint: Your redis should support lua script.

### LockLua
```lua
if (redis.call('EXISTS', KEYS[1]) == 0) then
    -- don't have lock
    redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
end
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 1) then
    -- reentry
    redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
end
return redis.call('PTTL', KEYS[1])
```

### UnLockLua
```lua
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 0) then
    -- not hold lock
    return -1
end
local counter = redis.call('HINCRBY', KEYS[1], ARGV[1], -1)
if (counter > 0) then
    -- update expire
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
else
    -- release lock
    redis.call('DEL', KEYS[1])
    return 1
end
return -1
```

### DelayExpireLua
```lua
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 1) then
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 1
end
return 0
```