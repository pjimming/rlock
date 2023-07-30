package rlock

import "time"

type Param struct {
	// redis distributedParam
	Addr     []string
	Password string

	// redis client timeout
	Timeout time.Duration

	// redis lock type
	Type string
}
