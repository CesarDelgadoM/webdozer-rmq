package database

import (
	"context"
	"log"
)

type IList interface {
	Pop(key string) string
	Add(key string, value interface{}) error
	Length(key string) int64
}

type List struct {
	*RedisPool
}

func NewList(client *RedisPool) *List {

	return &List{
		client,
	}
}

func (l *List) Pop(key string) string {

	value, err := l.client.LPop(context.Background(), key).Result()
	if err != nil {
		log.Println("Failed operation Pop:", err)
	}

	return value
}

func (l *List) Add(key string, value interface{}) error {
	return nil
}

func (l *List) Length(key string) int64 {

	size, err := l.client.LLen(context.Background(), key).Result()
	if err != nil {
		log.Println("Failed to calculate list length")
	}

	return size
}
