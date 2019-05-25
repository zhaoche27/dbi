package redis

import (
	"fmt"
	"time"

	"github.com/zhaoche27/dbi/dlock"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

var (
	keyEmptry   = fmt.Errorf("key is emptry")
	lockTimeout = fmt.Errorf("Get lock timeout")
)

// Distributed lock is Redis lock adapter.
type Lock struct {
	p *redis.Pool   // redis connection pool
	m string        // module
	k string        // key
	e time.Duration // expire
	v string        // value, default uuid
}

func NewLock(p *redis.Pool, m string) dlock.Lock {
	v := uuid.NewV4().String()
	return &Lock{p: p, m: m, v: v}
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

func (lock *Lock) TryLockAWaitInterval(key string, expire time.Duration,
	wait time.Duration, interval time.Duration) (bool, error) {
	lock.k = key
	lock.e = expire
	for stop := time.Now().Add(wait); ; {
		b, err := lock.nxSet()
		if b {
			return b, err
		}
		time.Sleep(interval)
		if time.Now().After(stop) {
			return false, lockTimeout
		}
	}
}

func (lock *Lock) Unlock() (bool, error) {
	if lock.k == "" {
		return false, keyEmptry
	}
	return lock.del()
}

func (lock *Lock) nxSet() (bool, error) {
	c := lock.p.Get()
	defer c.Close()
	if c.Err() != nil {
		return false, c.Err()
	}
	key := fmt.Sprintf("%s.%s", lock.m, lock.k)
	v, err := redis.Bool(c.Do("SET", key, lock.v, "NX", "EX", lock.e.Seconds))
	return v, err
}

func (lock *Lock) del() (bool, error) {
	c := lock.p.Get()
	defer c.Close()
	if c.Err() != nil {
		return false, c.Err()
	}
	luaScript := `
	if redis.call('get', KEYS[1]) == ARGV[1] 
    then 
        return redis.call('del', KEYS[1]) 
    else 
        return 0 
	end
	`
	key := fmt.Sprintf("%s.%s", lock.m, lock.k)
	v, err := redis.Bool(c.Do("EVAL", luaScript, 1, key, lock.v))
	return v, err
}
