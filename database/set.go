package database

import (
	"context"
	"log"
)

type ISet interface {
	Add(key string, value interface{})
	Remove(key string, value interface{})
	IsMember(key, value string) bool
	Length(key string) int64
	GetValues(key string) []string
}

type Set struct {
	ctx context.Context
	*RedisPool
}

func NewSet(redis *RedisPool) *Set {

	return &Set{
		context.Background(),
		redis,
	}
}

func (s *Set) Add(key string, value interface{}) {

	_, err := s.client.SAdd(s.ctx, key, value).Result()
	if err != nil {
		log.Println("Failed opertation SAdd:", err)
	}
}

func (s *Set) Remove(key string, value interface{}) {

	_, err := s.client.SRem(s.ctx, key, value).Result()
	if err != nil {
		log.Println("Failed operation SRem:", err)
	}
}

func (s *Set) IsMember(key, value string) bool {

	isMember, err := s.client.SIsMember(s.ctx, key, value).Result()
	if err != nil {
		log.Println("failed operation SIsMember:", err)
	}

	return isMember
}

func (s *Set) Length(key string) int64 {

	length, err := s.client.SCard(s.ctx, key).Result()
	if err != nil {
		log.Println("Failed operation SCard:", err)
	}

	return length
}

func (s *Set) GetValues(key string) []string {

	keys, err := s.client.SDiff(s.ctx, key).Result()
	if err != nil {
		log.Println("Failed operation SDiff:", err)
	}

	return keys
}
