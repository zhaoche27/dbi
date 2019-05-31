package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/zhaoche27/dbi/dlock"
)

// Distributed lock is Redis lock adapter.
type Lock struct {
	dlock.BaseLock
	p *redis.Pool   // redis connection pool
	m string        // module
	k string        // key
	e time.Duration // expire
	v string        // value, default uuid
	l bool
}

func NewLock(p *redis.Pool, m string) *Lock {
	v := uuid.NewV4().String()
	lock := &Lock{p: p, m: m, v: v}
	lock.Locker = lock
	return lock
}

func (lock *Lock) LockAWaitInterval(key string, expire time.Duration,
	wait time.Duration, interval time.Duration) (bool, error) {
	lock.k = key
	lock.e = expire
	lock.l = false
	for stop := time.Now().Add(wait); ; {
		b, err := lock.nxSet()
		if b {
			lock.l = true
			return b, err
		}
		time.Sleep(interval)
		if time.Now().After(stop) {
			return false, dlock.WaitTimeout
		}
	}
}

func (lock *Lock) Unlock() (bool, error) {
	if !lock.l {
		return false, nil
	}
	if lock.k == "" {
		return false, dlock.KeyEmptry
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
	v, err := redis.String(c.Do("SET", key, lock.v, "NX", "EX", lock.e.Seconds()))
	if v == "OK" {
		return true, err
	}
	return false, err
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
