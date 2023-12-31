package rlock

import "time"

const (
	// defaultBlockWaitingTime default block wait time to acquire lock
	defaultBlockWaitingTime = 60 * time.Second
	// defaultExpireTime default redis key expire time
	defaultExpireTime = 30 * time.Second
	// defaultWatchdogSwitch watchdog switch default false
	defaultWatchdogSwitch = false
)

// Lua Scripts.
// We use lua scripts for redis to ensure atomicity.
const (
	lockLua = `
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
`

	unLockLua = `
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
`

	delayExpireLua = `
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 1) then
	-- hold lock
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 1
end
return 0
`
)
