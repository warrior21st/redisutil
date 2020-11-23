package redisutil

import (
	"github.com/go-redis/redis"
)

//获取一个新redis连接
func GetNewClient(addr string, pwd string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
}
