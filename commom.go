package rlock

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// common variable
var (
	once sync.Once
	rc   *redis.Client
)

// generateToken 生成token
func generateToken() string {
	return fmt.Sprintf("%s_%s", getCurrentProcessID(), getCurrentGoroutineID())
}

// getCurrentProcessID 获取当前进程ID
func getCurrentProcessID() string {
	return strconv.Itoa(os.Getpid())
}

// getCurrentGoroutineID 获取当前的协程ID
func getCurrentGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
}
