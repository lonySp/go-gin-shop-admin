package models

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var redisCoretxt = context.Background()
var (
	RedisDb *redis.Client
)

func init() {
	fmt.Println("redis init.............")
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisDb.Ping(redisCoretxt).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
	} else {
		fmt.Println("Successfully connected to Redis.")
	}
}
