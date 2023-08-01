package utils

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// GenerateToken 生成token
func GenerateToken() string {
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
