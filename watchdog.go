package rlock

import (
	"context"
	"sync/atomic"
	"time"
)

func (l *RLock) watchdog() {
	// if watchdog switch off
	if !l.WatchdogSwitch() {
		return
	}

	// ensure the watchdog is stopping
	for !atomic.CompareAndSwapInt32(&l.runningDog, 0, 1) {
	}

	// run watchdog
	var ctx context.Context
	ctx, l.stopDog = context.WithCancel(l.ctx)
	go func() {
		defer func() {
			atomic.StoreInt32(&l.runningDog, 0)
		}()
		l.runWatchdog(ctx)
	}()
}

func (l *RLock) runWatchdog(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(l.expireSeconds*1000/3) * time.Millisecond)
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
