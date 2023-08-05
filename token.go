package rlock

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// generateToken 生成token
func generateToken() string {
	return fmt.Sprintf("%s_%s_%s", getCurrentProcessID(), getCurrentGoroutineID(), generateRandomString(6))
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

// generateRandomString 生成长度为length的随机字符串
func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	ret := ""
	for i := 0; i < length; i++ {
		ret += string(charset[rand.Intn(len(charset))])
	}
	return ret
}
