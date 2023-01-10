package redisutil

import (
	"strings"

	"github.com/go-redis/redis"
)

func Enqueue(rc *redis.Client, queueName string, val string) {
	rc.LPush(queueName, val)
}

func Dequeue(rc *redis.Client, queueName string) string {
	b, err := rc.RPopLPush(queueName, getDequeuedBackupKey(queueName)).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	if b == "" {
		return ""
	}

	return b
}

func RestoreDequeuedBackup(client *redis.Client, queueName string) {
	backupKey := getDequeuedBackupKey(queueName)
	for {
		err := client.RPopLPush(backupKey, queueName).Err()
		if err != nil {
			if err == redis.Nil {
				break
			} else {
				panic(err)
			}
		}
	}
}

func DelLastDequeuedBackup(client *redis.Client, queueName string) {
	client.LTrim(getDequeuedBackupKey(queueName), 1, -1)
}

func getDequeuedBackupKey(queueName string) string {
	return strings.Join([]string{queueName, "_dequeued_backup"}, "")
}
