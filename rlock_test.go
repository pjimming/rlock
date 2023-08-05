package rlock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var op = RedisClientOptions{
	Addr:     "127.0.0.1:6379",
	Password: "",
}

func TestNewRLock(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "")
	ast.NotNil(l)
}

func TestLock(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "")
	ast.NotNil(l)

	ttl := l.Lock()
	ast.Equal(int64(0), ttl)

	l.UnLock()
}

func TestTryLock(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "")
	ast.NotNil(l)

	ttl := l.TryLock()
	t.Log("ttl:", ttl)
	ast.Equal(int64(0), ttl)

	l.UnLock()
}

func TestLockTwice(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "").SetWatchdogSwitch(true).SetExpireTime(3 * time.Second)
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

	l := NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.TryLock()
	ast.Equal(int64(0), ttl)

	ttl2 := l.TryLock()
	ast.Equal(int64(0), ttl2)
}

func TestLockAndUnLock(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.Lock()
	ast.Equal(int64(0), ttl)

	res := l.UnLock()
	ast.Equal(int64(1), res)
}

func TestTryLockAndUnLock(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ttl := l.TryLock()
	ast.Equal(int64(0), ttl)

	res := l.UnLock()
	ast.Equal(int64(1), res)
}

func TestReentry(t *testing.T) {
	ast := assert.New(t)

	l := NewRLock(op, "")
	ast.NotNil(l)

	t.Log(l.Key(), l.Token())

	ast.Equal(int64(0), l.Lock())
	ast.Equal(int64(0), l.Lock())

	ast.Equal(int64(0), l.UnLock())
	ast.Equal(int64(1), l.UnLock())
}

func TestBlocking(t *testing.T) {
	ast := assert.New(t)
	key := "22229999"

	l1 := NewRLock(op, key).
		SetToken(key + "111").
		SetExpireTime(5 * time.Second)

	l2 := NewRLock(op, key).
		SetToken(key + "222").
		SetBlockWaitingSecond(20 * time.Second)

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
	key := "22229999000"

	l1 := NewRLock(op, key).
		SetToken(key + "111").
		SetExpireTime(5 * time.Second).
		SetWatchdogSwitch(true)

	l2 := NewRLock(op, key).
		SetToken(key + "222").
		SetBlockWaitingSecond(20 * time.Second)

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

func TestRedLock(t *testing.T) {
	redLock, err := NewRedLock([]RedisClientOptions{
		{Addr: "127.0.0.1:7001", Password: ""},
		{Addr: "127.0.0.1:7002", Password: ""},
		{Addr: "127.0.0.1:7003", Password: ""},
		{Addr: "127.0.0.1:7004", Password: ""},
		{Addr: "127.0.0.1:7005", Password: ""},
	}, "1234567_key", 30*time.Second)

	if err != nil {
		t.Log(err)
		return
	}

	t.Log(redLock.TryLock())
	redLock.UnLock()
}
