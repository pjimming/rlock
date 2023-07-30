# Redis Lock

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Usage
```shell
go get -u github.com/pjimming/rlock
```

## Quick Start

```go
package main

import (
	"time"

	"github.com/pjimming/rlock"
)

func main() {
	l, err := rlock.NewRedisLock(rlock.Param{
		Addr:     []string{"127.0.0.1:6379"},
		Password: "",
		Timeout:  100,
		Type:     "",
	})
	if err != nil {
		return
	}

	key := "key"
	value := "value"

	// Try Lock
	_, err = l.TryLock(key, value, 10*time.Second)
	if err != nil {
		return
	}

	// Release Lock
	_, err = l.ReleaseLock(key, value)
	if err != nil {
		return
	}
}
```

## [Lua Scripts](./lua.md)
> Hint: Your redis should support lua script.