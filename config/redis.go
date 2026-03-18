package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {

	if os.Getenv("REDISHOST") == "" {
		log.Fatal("REDISHOST is empty")
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDISHOST"),
		Username:     os.Getenv("REDISUSER"),
		Password:     os.Getenv("REDISPASS"),
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
}
