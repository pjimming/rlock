package rlock

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

var reentryParam = Param{
	Addr:     []string{"127.0.0.1:6379"},
	Password: "",
	Timeout:  100,
	Type:     TypeReentry,
}

func TestRLockReentry_TryLockAndReleaseLock(t *testing.T) {
	key := getRandomString()
	value := getRandomString()

	ast := assert.New(t)

	l, err := NewRedisLock(reentryParam)
	ast.Nil(err)
	ast.NotNil(l)

	isLock, err := l.TryLock(key, value, expireTime)
	ast.Nil(err)
	ast.Equal(true, isLock)

	ast.Equal(true, l.HoldLock(key, value))
	ast.Equal(false, l.HoldLock(key, value+"111"))

	isRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(true, isRelease)
}

func TestRLockReentry_TryLockTwice(t *testing.T) {
	key := getRandomString()
	value := getRandomString()

	ast := assert.New(t)

	l, err := NewRedisLock(reentryParam)
	ast.Nil(err)
	ast.NotNil(l)

	ast.Equal(false, l.HoldLock(key, value))

	firstLock, err := l.TryLock(key, value, expireTime)
	ast.Nil(err)
	ast.Equal(true, firstLock)

	ast.Equal(true, l.HoldLock(key, value))

	secondLock, err := l.TryLock(key, value, expireTime)
	ast.Nil(err)
	ast.Equal(true, secondLock)

	ast.Equal(true, l.HoldLock(key, value))

	firstRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(false, firstRelease)

	ast.Equal(true, l.HoldLock(key, value))

	secondRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(true, secondRelease)

	ast.Equal(false, l.HoldLock(key, value))
}

func TestRLockReentry_ReleaseEmpty(t *testing.T) {
	ast := assert.New(t)
	l, err := NewRedisLock(reentryParam)
	ast.Nil(err)
	ast.NotNil(l)

	release, err := l.ReleaseLock(getRandomString(), getRandomString())
	ast.Nil(err)
	ast.Equal(false, release)
}

func TestRLockReentry_Goroutine(t *testing.T) {
	key := getRandomString()
	ast := assert.New(t)

	var wg sync.WaitGroup
	var mu sync.Mutex

	successCnt, failCnt := 0, 0
	for i := 0; i < 10; i++ {
		vv := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			getLock := false

			l, err := NewRedisLock(reentryParam)
			ast.Nil(err)
			ast.NotNil(l)

			value := strconv.Itoa(vv)

			for i := 0; i < 10; i++ {
				isLock, err := l.TryLock(key, value, expireTime)
				ast.Nil(err)

				mu.Lock()
				if isLock {
					successCnt++
					getLock = true
				} else {
					failCnt++
				}
				mu.Unlock()
			}

			dlog.Infof("goroutine %d get lock %v", vv, getLock)

			for i := 0; i < 10; i++ {
				isRelease, err := l.ReleaseLock(key, value)
				ast.Nil(err)

				if isRelease {
					dlog.Infof("goroutine %d release lock %v", vv, isRelease)
				}
			}
		}()
	}
	wg.Wait()

	ast.Equal(10, successCnt)
	ast.Equal(90, failCnt)
}
