package rlock

import (
	"context"
	"sync/atomic"
	"time"
)

// watchdog new a watchdog to delay redis lock's expire time.
func (l *RLock) watchdog() {
	// if watchdog switch off
	if !l.WatchdogSwitch() {
		return
	}

	// ensure the watchdog is stopping
	for !atomic.CompareAndSwapInt32(&l.runningDog, 0, 1) {
	}

	// run watchdog, set stopDog func
	var ctx context.Context
	ctx, l.stopDog = context.WithCancel(l.ctx)
	go func() {
		defer func() {
			atomic.StoreInt32(&l.runningDog, 0)
		}()
		l.runWatchdog(ctx)
	}()
}

// runWatchdog start to expire redis lock.
func (l *RLock) runWatchdog(ctx context.Context) {
	ticker := time.NewTicker(l.expireTime / 3)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if _, err := l.delayExpire(); err != nil {
			return
		}
	}
}
