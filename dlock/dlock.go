package dlock

import (
	"fmt"
	"time"
)

const (
	DefaultInteval = 30 * time.Microsecond
	DefaultWait    = 3 * time.Hour
)

var (
	KeyEmptry   = fmt.Errorf("key is emptry")
	WaitTimeout = fmt.Errorf("Get lock wait timeout")
)

// Distributed lock interface
type Locker interface {
	Lock(key string, expire time.Duration) error
	TryLock(key string, expire time.Duration) (bool, error)
	TryLockAwait(key string, expire time.Duration, wait time.Duration) (bool, error)
	TryLockAWaitInterval(key string, expire time.Duration, wait time.Duration, interval time.Duration) (bool, error)
	LockAWait(key string, expire time.Duration, wait time.Duration) (bool, error)
	LockAWaitInterval(key string, expire time.Duration, wait time.Duration, interval time.Duration) (bool, error)
	Unlock() (bool, error)
}

type BaseLock struct {
	Locker
}

func (lock *BaseLock) Lock(key string, expire time.Duration) error {
	_, err := lock.LockAWait(key, expire, DefaultWait)
	return err
}
func (lock *BaseLock) LockAWait(key string, expire time.Duration, wait time.Duration) (bool, error) {
	_, err := lock.Locker.LockAWaitInterval(key, expire, wait, DefaultInteval)
	return false, err
}
func (lock *BaseLock) TryLock(key string, expire time.Duration) (bool, error) {
	return lock.TryLockAwait(key, expire, 0)
}
func (lock *BaseLock) TryLockAwait(key string, expire time.Duration, wait time.Duration) (bool, error) {
	return lock.TryLockAWaitInterval(key, expire, wait, DefaultInteval)
}
func (lock *BaseLock) TryLockAWaitInterval(key string, expire time.Duration,
	wait time.Duration, interval time.Duration) (bool, error) {
	b, err := lock.Locker.LockAWaitInterval(key, expire, wait, interval)
	if err == WaitTimeout {
		return b, nil
	}
	return b, err
}
