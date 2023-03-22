package redisutil

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

func Enqueue(rc *redis.Conn, queueName string, val string) {

	rc.LPush(context.Background(), queueName, val)
}

func Dequeue(rc *redis.Conn, queueName string) string {
	b, err := rc.RPopLPush(context.Background(), queueName, getDequeuedBackupKey(queueName)).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	if b == "" {
		return ""
	}

	return b
}

func RestoreDequeuedBackup(client *redis.Conn, queueName string) {
	backupKey := getDequeuedBackupKey(queueName)
	for {
		err := client.RPopLPush(context.Background(), backupKey, queueName).Err()
		if err != nil {
			if err == redis.Nil {
				break
			} else {
				panic(err)
			}
		}
	}
}

func DelLastDequeuedBackup(client *redis.Conn, queueName string) {
	client.LTrim(context.Background(), getDequeuedBackupKey(queueName), 1, -1)
}

func getDequeuedBackupKey(queueName string) string {
	return strings.Join([]string{queueName, "_dequeued_backup"}, "")
}
