package database

import (
	"context"
	"log"
	"strings"
)

type IHash interface {
	GetValues(key string) []string
	GetKeys(hash string) []string
}

type Hash struct {
	name string
	ctx  context.Context
	*RedisPool
}

func NewHash(name string, redis *RedisPool) *Hash {

	return &Hash{
		name,
		context.Background(),
		redis,
	}
}

func (h *Hash) GetValues(key string) []string {

	value, err := h.client.HGet(h.ctx, h.name, key).Result()
	if err != nil {
		log.Println("Failed operation HGet:", err)
	}

	return strings.Split(value, ",")
}

func (h *Hash) GetKeys(hash string) []string {

	keys, err := h.client.HKeys(h.ctx, hash).Result()
	if err != nil {
		log.Println("Failed operation HKeys:", err)
	}

	return keys
}
