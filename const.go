package rlock

// Lua Scripts
const (
	commonLockLua = `
local lockKey = KEYS[1]
local lockToken = ARGV[1]
local expireTime = tonumber(ARGV[2])

SET KEYS[1] ARGV[1] NX EX expireTime 
`

	reentryLockLua = `
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
	reentryUnlockLua = `
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

	expireLua = ``
)
