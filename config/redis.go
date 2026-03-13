package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {

	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISHOST"),
		Username: os.Getenv("REDISUSER"),
		Password: os.Getenv("REDISPASS"),
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
}
