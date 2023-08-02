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

func TestNewRLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)
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

func TestLockTwice(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "").SetWatchdogSwitch(true).SetExpireSeconds(3)
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.Lock()
	ast.Equal(int64(0), ttl)

	ttl2 := l.Lock()
	ast.Equal(int64(0), ttl2)

	ast.Equal(int64(0), l.UnLock())
	ast.Equal(int64(1), l.UnLock())

	time.Sleep(10 * time.Second)
}

func TestTryLockTwice(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.TryLock()
	ast.Equal(int64(0), ttl)

	ttl2 := l.TryLock()
	ast.Equal(int64(0), ttl2)
}

func TestLockAndUnLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.Lock()
	ast.Equal(int64(0), ttl)

	res := l.UnLock()
	ast.Equal(int64(1), res)
}

func TestTryLockAndUnLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.TryLock()
	ast.Equal(int64(0), ttl)

	res := l.UnLock()
	ast.Equal(int64(1), res)
}

func TestReentry(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ast.Equal(int64(0), l.Lock())
	ast.Equal(int64(0), l.Lock())

	ast.Equal(int64(0), l.UnLock())
	ast.Equal(int64(1), l.UnLock())
}

func TestBlocking(t *testing.T) {
	ast := assert.New(t)
	key := utils.GenerateRandomString(8)

	l1 := rlock.NewRLock(op, key).SetToken(key + "111").SetExpireSeconds(5)
	l2 := rlock.NewRLock(op, key).SetToken(key + "222").SetBlockWaitingSecond(20)

	t.Log("l1:", l1.Key(), l1.Token())
	t.Log("l2:", l2.Key(), l2.Token())

	ast.Equal(int64(0), l1.Lock())

	start := time.Now()
	ast.Less(int64(0), l2.TryLock())
	t.Log("l2 TryLock cost:", time.Now().Sub(start).String())

	t.Log("l1 start sleep...", time.Now().Sub(start).String())
	time.Sleep(time.Second * 4)
	ast.Equal(int64(1), l1.UnLock())
	t.Log("l1 unlock", time.Now().Sub(start).String())

	ast.Equal(int64(0), l2.Lock())
	t.Log("l2 Lock cost:", time.Now().Sub(start).String())

	ast.Equal(int64(1), l2.UnLock())
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
