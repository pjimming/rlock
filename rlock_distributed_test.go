package rlock

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const expireTime = 10 * time.Second

var distributedParam = Param{
	Addr:     []string{"127.0.0.1:6379"},
	Password: "",
	Timeout:  100,
	Type:     TypeDistributed,
}

func TestTryLockAndReleaseLock(t *testing.T) {
	key := getRandomString()
	value := getRandomString()

	ast := assert.New(t)
	l, err := NewRedisLock(distributedParam)
	ast.Nil(err)
	ast.NotNil(l)

	// TryLock
	isLock, err := l.TryLock(key, value, expireTime)
	ast.Nil(err)
	ast.Equal(true, isLock)

	// IsOwner
	ast.Equal(true, l.HoldLock(key, value))

	// ReleaseLock
	isRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(true, isRelease)
}

func TestTryLockTwice(t *testing.T) {
	ast := assert.New(t)

	key := getRandomString()

	l, err := NewRedisLock(distributedParam)
	ast.Nil(err)
	ast.NotNil(l)

	firstLock, err := l.TryLock(key, getRandomString(), expireTime)
	ast.Nil(err)
	ast.Equal(true, firstLock)

	secondLock, err := l.TryLock(key, getRandomString(), expireTime)
	ast.Nil(err)
	ast.Equal(false, secondLock)
}

func TestReleaseLockEmpty(t *testing.T) {
	key := getRandomString()

	ast := assert.New(t)

	l, err := NewRedisLock(distributedParam)
	ast.Nil(err)
	ast.NotNil(l)

	isRelease, err := l.ReleaseLock(key, getRandomString())
	ast.Nil(err)
	ast.Equal(false, isRelease)
}

func TestReleaseTwice(t *testing.T) {
	key := getRandomString()
	value := getRandomString()

	ast := assert.New(t)

	l, err := NewRedisLock(distributedParam)
	ast.Nil(err)
	ast.NotNil(l)

	isLock, err := l.TryLock(key, value, expireTime)
	ast.Nil(err)
	ast.Equal(true, isLock)

	ast.Equal(true, l.HoldLock(key, value))

	firstRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(true, firstRelease)

	ast.Equal(false, l.HoldLock(key, value))

	secondRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(false, secondRelease)
}

func TestReleaseOther(t *testing.T) {
	key := getRandomString()
	value := getRandomString()

	ast := assert.New(t)

	l, err := NewRedisLock(distributedParam)
	ast.Nil(err)
	ast.NotNil(l)

	isLock, err := l.TryLock(key, value, expireTime)
	ast.Nil(err)
	ast.Equal(true, isLock)

	otherRelease, err := l.ReleaseLock(key, value+"111")
	ast.Nil(err)
	ast.Equal(false, otherRelease)

	ast.Equal(true, l.HoldLock(key, value))

	isRelease, err := l.ReleaseLock(key, value)
	ast.Nil(err)
	ast.Equal(true, isRelease)

	ast.Equal(false, l.HoldLock(key, value))
}

func TestGoroutineTryLock(t *testing.T) {
	key := getRandomString()

	ast := assert.New(t)

	var wg sync.WaitGroup
	var mu sync.Mutex

	successCnt, failCnt := 0, 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l, err := NewRedisLock(distributedParam)
			ast.Nil(err)
			isLock, err := l.TryLock(key, getRandomString(), expireTime)
			ast.Nil(err)
			fmt.Printf("goroutine %d get lock %v\n", i, isLock)

			mu.Lock()
			if isLock {
				successCnt++
			} else {
				failCnt++
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	ast.Equal(1, successCnt)
	ast.Equal(9, failCnt)
}

func getRandomString() string {
	rand.Seed(time.Now().UnixNano())
	var ret string
	for i := 0; i < 10; i++ {
		randomNumber := rand.Intn(10) // 生成 0 到 9 的随机整数
		ret += strconv.Itoa(randomNumber)
	}
	return ret
}
