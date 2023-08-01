package test

import (
	"testing"

	"github.com/pjimming/rlock"

	"github.com/stretchr/testify/assert"
)

var op = rlock.RedisClientOptions{
	Addr:     "127.0.0.1:6379",
	Password: "",
}

func TestNewRLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op)
	ast.NotNil(l)
}
