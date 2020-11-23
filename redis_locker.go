package redisutil

import (
	"time"

	"github.com/go-redis/redis"
)

const _lockHoldersHashSet = "all_lock_holders"

//redis 锁
type RedisLock struct {
	Key           string //锁key
	Holder        string //锁拥有者
	ExpireSeconds int64  //锁过期时间（秒）
}

//获取锁（永不过期）
func AcquireLock(client *redis.Client, key string, holder string) (*RedisLock, error) {
	return AcquireLockWithExpire(client, key, holder, 0)
}

//获取锁（带过期时间）
func AcquireLockWithExpire(client *redis.Client, key string, holder string, expireSeconds int64) (*RedisLock, error) {
	var (
		lock *RedisLock
		err  error
		b    bool  = false
		temp int32 = 0
	)
	for lock == nil {
		b, err = client.SetNX(key, "1", 0).Result()
		if err != nil {
			return lock, err
		}
		if b {
			client.HSet(_lockHoldersHashSet, key, holder)
			if expireSeconds > 0 {
				client.Expire(key, time.Second*time.Duration(expireSeconds))
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
func ReleaseLock(client *redis.Client, key string, holder string) (bool, error) {
	b := false
	var err error
	currHolder, err := client.HGet(_lockHoldersHashSet, key).Result()
	if err != nil {
		return b, err
	}
	if currHolder == holder {
		client.Del(key)
		client.HDel(_lockHoldersHashSet, key)
		b = true
	}

	return b, nil
}

//释放
func (lock *RedisLock) ReleaseSelf(client *redis.Client) (bool, error) {
	b := false
	var err error
	currHolder, err := client.HGet(_lockHoldersHashSet, lock.Key).Result()
	if err != nil {
		return b, err
	}
	if currHolder == lock.Holder {
		client.Del(lock.Key)
		client.HDel(_lockHoldersHashSet, lock.Key)
		b = true
	}

	return b, nil
}

//清除指定holder的拥有的所有锁
func ClearLocks(client *redis.Client, holder string) error {
	holderLocks, err := client.HGetAll(_lockHoldersHashSet).Result()
	if err != nil {
		return err
	}
	for key, val := range holderLocks {
		if val == holder {
			client.Del(key)
			client.HDel(_lockHoldersHashSet, key)
		}
	}

	return nil
}
