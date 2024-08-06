package cache

import (
	"context"
	"log"
	"os"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RDB * redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error!")
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

		addr := host + ":" + port

		RDB = redis.NewClient( & redis.Options {
		Addr: addr,
		Password: password,
		DB: 0,
	})

		_,
	err = RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Error:", err)
	}
}
