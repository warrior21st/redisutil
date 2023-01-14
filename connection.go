package redisutil

import (
	"github.com/go-redis/redis"
)

func GetNewClient(connStr string) *redis.Client {
	redisOpt, err := redis.ParseURL(connStr)
	if err != nil {
		panic(err)
	}

	// redisOpt.DialTimeout = time.Hour * 24 * 365 * 100
	// redisOpt.IdleTimeout = time.Hour * 24 * 365 * 100
	// redisOpt.PoolTimeout = time.Hour * 24 * 365 * 100
	// redisOpt.ReadTimeout = time.Hour * 24 * 365 * 100
	// redisOpt.WriteTimeout = time.Hour * 24 * 365 * 100

	// log.Printf("DialTimeout: %d", redisOpt.DialTimeout/time.Millisecond)
	// log.Printf("IdleTimeout: %d", redisOpt.IdleTimeout/time.Millisecond)
	// log.Printf("PoolTimeout: %d", redisOpt.PoolTimeout/time.Millisecond)
	// log.Printf("ReadTimeout: %d", redisOpt.ReadTimeout/time.Millisecond)
	// log.Printf("WriteTimeout: %d", redisOpt.WriteTimeout/time.Millisecond)

	return redis.NewClient(redisOpt)
}
