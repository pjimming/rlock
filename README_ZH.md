# RLock

由Golang实现的基于Redis的分布式锁。

- [English](./README.md)

---

## 状态
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 使用
```shell
go get -u github.com/pjimming/rlock
```

## 实现功能
- 互斥性：Redis分布式锁可以保证同一时刻只有一个客户端可以获得锁，实现线程之间的互斥。
- 安全性：Redis分布式锁采用原子操作，可以保证并发情况下锁的安全性，避免数据竞争、死锁等问题。
- 锁超时：为了避免某个客户端获取锁后失败而导致死锁，Redis分布式锁可以设置锁超时时间，超过超时时间会自动释放锁。
- 可重入性：Redis分布式锁可以支持同一个客户端多次获取同一个锁，避免嵌套调用时出现死锁。
- 高性能：Redis是一个内存数据库，具有很高的读写性能，可以实现快速的加锁和解锁操作。
- 原子性：Redis分布式锁的加锁和解锁操作使用原子命令，可以保证操作的原子性，避免并发下的竞争问题。

## 快速开始
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

## Lua 脚本
> Hint: Your redis should support lua script.

### 加锁Lua
```lua
if (redis.call('EXISTS', KEYS[1]) == 0) then
    -- 锁未被占有
    redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
end
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 1) then
    -- 可重入
    redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
end
return redis.call('PTTL', KEYS[1])
```

### 解锁Lua
```lua
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 0) then
    -- 未持有锁
    return -1
end
local counter = redis.call('HINCRBY', KEYS[1], ARGV[1], -1)
if (counter > 0) then
    -- 更新过期时间
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 0
else
    -- 释放锁
    redis.call('DEL', KEYS[1])
    return 1
end
return -1
```

### 续命
```lua
if (redis.call('HEXISTS', KEYS[1], ARGV[1]) == 1) then
    redis.call('PEXPIRE', KEYS[1], tonumber(ARGV[2]))
    return 1
end
return 0
```