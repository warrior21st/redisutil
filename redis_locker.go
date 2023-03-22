package redisutil

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const _lockHoldersHashSet = "all_lock_holders"

//redis 锁
type RedisLock struct {
	Key           string //锁key
	Holder        string //锁拥有者
	ExpireSeconds int64  //锁过期时间（秒）
}

//获取锁（永不过期）
func AcquireLock(client *redis.Conn, key string, holder string) (*RedisLock, error) {
	return AcquireLockWithExpire(client, key, holder, 0)
}

//获取锁（带过期时间）
func AcquireLockWithExpire(conn *redis.Conn, key string, holder string, expireSeconds int64) (*RedisLock, error) {
	var (
		lock *RedisLock
		err  error
		b    bool  = false
		temp int32 = 0
	)
	for lock == nil {
		b, err = conn.SetNX(context.Background(), key, "1", 0).Result()
		if err != nil {
			return lock, err
		}
		if b {
			conn.HSet(context.Background(), _lockHoldersHashSet, key, holder)
			if expireSeconds > 0 {
				conn.Expire(context.Background(), key, time.Second*time.Duration(expireSeconds))
			}
			lock = &RedisLock{
				Key:           key,
				Holder:        holder,
				ExpireSeconds: expireSeconds,
			}
		} else {
			time.Sleep(time.Millisecond * 10)
			temp += 10
			//最多持续争抢2分钟
			if temp > 1000*120 {
				break
			}
		}
	}
	return lock, nil
}

//释放锁
func ReleaseLock(client *redis.Conn, key string, holder string) (bool, error) {
	b := false
	var err error
	currHolder, err := client.HGet(context.Background(), _lockHoldersHashSet, key).Result()
	if err != nil {
		return b, err
	}
	if currHolder == holder {
		client.Del(context.Background(), key)
		client.HDel(context.Background(), _lockHoldersHashSet, key)
		b = true
	}

	return b, nil
}

//释放
func (lock *RedisLock) ReleaseSelf(client *redis.Conn) (bool, error) {
	b := false
	var err error
	currHolder, err := client.HGet(context.Background(), _lockHoldersHashSet, lock.Key).Result()
	if err != nil {
		return b, err
	}
	if currHolder == lock.Holder {
		client.Del(context.Background(), lock.Key)
		client.HDel(context.Background(), _lockHoldersHashSet, lock.Key)
		b = true
	}

	return b, nil
}

//清除指定holder的拥有的所有锁
func ClearLocks(client *redis.Conn, holder string) error {
	holderLocks, err := client.HGetAll(context.Background(), _lockHoldersHashSet).Result()
	if err != nil {
		return err
	}
	for key, val := range holderLocks {
		if val == holder {
			client.Del(context.Background(), key)
			client.HDel(context.Background(), _lockHoldersHashSet, key)
		}
	}

	return nil
}
