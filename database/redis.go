package database

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisPool struct {
	client *redis.Client
}

func NewRedisPool(addr string, procs int) *RedisPool {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		PoolSize: procs,
	})

	log.Println("redis:PING... ")

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed connect to redis")
	}

	log.Println("redis:" + pong + "!")

	return &RedisPool{
		client: client,
	}
}

func (rp *RedisPool) GetKeys(pattern string) []string {

	keys, err := rp.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		log.Fatal("Failed to catch keys")
	}

	return keys
}

func (rp *RedisPool) DelKey(key string) int64 {

	id, err := rp.client.Del(context.Background(), key).Result()
	if err != nil {
		log.Fatal("Failed to delete key")
	}

	return id
}

func (rp *RedisPool) Close() {

	err := rp.client.Close()
	if err != nil {
		log.Fatal("Failed to close redis connection")
	}
}
