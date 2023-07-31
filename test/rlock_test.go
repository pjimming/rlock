package test

import (
	"github.com/pjimming/rlock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var op = rlock.SingleClientOptions{
	Addr:     "127.0.0.1:6379",
	Password: "",
}

func TestNewRLock(t *testing.T) {
	ast := assert.New(t)

	l := rlock.NewRLock(op)
	ast.NotNil(l)
}
