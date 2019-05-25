package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestLock(t *testing.T) {
	lock := NewLock(pool(), "test")
	fmt.Println("-------------", lock)
	b, err := lock.TryLockAWaitInterval("key1", 20*time.Second, 20*time.Second, 5*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("0000", b)
	defer lock.Unlock()
	time.Sleep(10 * time.Second)
}

func pool() *redis.Pool {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", "127.0.0.1:6666")
		return c, err
	}
	// initialize a new pool
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}
