package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/connection"
	"os"
)

var world *World

func Init() {
	redisServer := os.Getenv("REDISTOGO_URL")
	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	redisConnection := connection.NewRedisConnection(redisServer)
	cache := cache.NewCache(redisConnection)
	w, err := newWorld(cache, 10)

	if err != nil {
		panic(err)
	}

	world = w
}
