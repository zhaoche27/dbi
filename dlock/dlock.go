package dlock

import (
	"time"
)

const (
	DefaultInteval = 30 * time.Microsecond
	DefaultWait    = 3 * time.Hour
)

// Distributed lock interface
type Lock interface {
	Lock(key string, expire time.Duration) error
	TryLock(key string, expire time.Duration) (bool, error)
	TryLockAwait(key string, expire time.Duration, wait time.Duration) (bool, error)
	TryLockAWaitInterval(key string, expire time.Duration, wait time.Duration, interval time.Duration) (bool, error)
	Unlock() (bool, error)
}
