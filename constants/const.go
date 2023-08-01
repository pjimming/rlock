package constants

const (
	ErrLockAcquiredByOthersStr = "lock acquire by others"
)

// Lua Scripts
const (
	LockLua = `
if redis.call("SET", KEYS[1], ARGV[1], "NX", "EX", tonumber(ARGV[2])) then
	return 1
else
	return 0
end
`

	UnLockLua = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else 
	return 0
end
`

	ReentryLockLua = `
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
	ReentryUnlockLua = `
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

	UpdateExpireLua = ``
)
