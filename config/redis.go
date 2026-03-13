package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {

	// conenction local
	// RDB = redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// 	DB:       0,
	// })

	// conenction cloud
	// RDB = redis.NewClient(&redis.Options{
	// 	// Addr: "localhost:6379",
	// 	Addr:     "redis-13734.crce185.ap-seast-1-1.ec2.cloud.redislabs.com:13734",
	// 	Username: "default",
	// 	Password: "LftGnlhE11Ex73HzjhVr4keSXUnTMFbp",
	// })

	// railway connection
	RDB = redis.NewClient(&redis.Options{
		// Addr: "localhost:6379",

		// redis://default:NxNrRWkkqUPdcGNkaUqJfhQLPEzqZaqA@redis.railway.internal:6379
		// redis://default:NxNrRWkkqUPdcGNkaUqJfhQLPEzqZaqA@mainline.proxy.rlwy.net:49575
		Addr:     "redis.railway.internal:6379",
		Username: "default",
		Password: "NxNrRWkkqUPdcGNkaUqJfhQLPEzqZaqA",
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
}
