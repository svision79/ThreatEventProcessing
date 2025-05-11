package repository

import (
	"github.com/redis/go-redis/v9"
	"os"
)

func NewCache() *redis.Client {
	add := os.Getenv("REDIS_ADDR")

	return redis.NewClient(&redis.Options{
		Addr:     add,
		Password: "",
		DB:       0,
	})
}
