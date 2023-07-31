package mistake

import (
	"errors"

	"github.com/pjimming/rlock/common"
)

func ErrLockAcquiredByOthers() error {
	return errors.New(common.ErrLockAcquiredByOthersStr)
}
