package redis_distributed_lock

import (
	"context"
	"errors"
	"sync"
	"testing"
)

// 阻塞分布式锁测试
func Test_blockingLock(t *testing.T) {
	// 请输入 redis 节点的地址和密码
	addr := "127.0.0.1:6379"
	passwd := ""
	client := NewClient("tcp", addr, passwd)
	lock1 := NewRedisLock("test_key", client, SetExpireSeconds(1))
	lock2 := NewRedisLock("test_key", client, ActiveBlockMode(), SetBlockWaitingSeconds(2))
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lock1.Lock(ctx); err != nil {
			t.Error(err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lock2.Lock(ctx); err != nil {
			t.Error(err)
			return
		}
	}()
	wg.Wait()
	t.Log("success")
}

// 非阻塞分布式锁
func Test_nonblockingLock(t *testing.T) {
	// 请输入 redis 节点的地址和密码
	addr := "127.0.0.1:6379"
	passwd := ""
	client := NewClient("tcp", addr, passwd)
	lock1 := NewRedisLock("test_key", client, SetExpireSeconds(1))
	lock2 := NewRedisLock("test_key", client)
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lock1.Lock(ctx); err != nil {
			t.Error(err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lock2.Lock(ctx); err == nil || !errors.Is(err, ErrLockAcquiredByOthers) {
			t.Errorf("got err: %v, expect: %v", err, ErrLockAcquiredByOthers)
			return
		}
	}()
	wg.Wait()
	t.Log("success")
}
