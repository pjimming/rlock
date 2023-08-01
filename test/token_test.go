package test

import (
	"testing"

	"github.com/pjimming/rlock/utils"
)

func TestGenerateToken(t *testing.T) {
	t.Log(utils.GenerateToken())
}
