package db

import (
	"time"

	"github.com/zhaoche27/dbi/dlock"
)

type Lock struct {
}

func (lock *Lock) Lock(key string, expire time.Duration) error {
	_, err := lock.TryLockAwait(key, expire, dlock.DefaultWait)
	return err
}
func (lock *Lock) TryLock(key string, expire time.Duration) (bool, error) {
	return lock.TryLockAwait(key, expire, 0)
}
func (lock *Lock) TryLockAwait(key string, expire time.Duration, wait time.Duration) (bool, error) {
	return lock.TryLockAWaitInterval(key, expire, wait, dlock.DefaultInteval)
}
func (lock *Lock) TryLockAWaitInterval(key string, expire time.Duration, wait time.Duration, interval time.Duration) (bool, error) {
	return false, nil
}
func (lock *Lock) Unlock() (bool, error) {
	return false, nil
}
