package mistake

import (
	"errors"

	"github.com/pjimming/rlock/constants"
)

func ErrLockAcquiredByOthers() error {
	return errors.New(constants.ErrLockAcquiredByOthersStr)
}
