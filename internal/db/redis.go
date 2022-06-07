package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisDB struct {
	ctx    context.Context
	client *redis.Client
}

func (r *RedisDB) Get(k string) (string, error) {
	val, err := r.client.Get(k).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("Key %s does not exist", k)
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisDB) Put(k, v string) error {
	err := r.client.Set(k, v, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisClient(addr string) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return &RedisDB{
		client: client,
		ctx:    ctx,
	}, nil
}
