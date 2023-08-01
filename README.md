# RLock

A distributed redis lock based by Golang.

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

## Lua Scripts
> Hint: Your redis should support lua script.

### LockLua
```lua
if (redis.call('EXISTS', KEYS[1]) == 0) then
    -- don't have lock
    redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
    redis.call('EXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
end
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 1) then
    -- reentry
    redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
    redis.call('EXPIRE', KEYS[1], tonumber(ARGV[2]))
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
    redis.call('EXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
else
    -- release lock
    redis.call('DEL', KEYS[1])
    return 1
end
return -1
```