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
	userCache := cache.NewUserCache(redisConnection)
	world = newWorld(userCache)
}
