package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/connection"
	"github.com/jpcummins/geochat/app/types"
	"os"
)

var world types.World

func Init() {
	redisServer := os.Getenv("REDISTOGO_URL")
	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	redisConnection := connection.NewRedisConnection(redisServer)
	cache := cache.NewCache(redisConnection)

	factory := &Factory{}
	w, err := factory.NewWorld(cache, 10)

	if err != nil {
		panic(err)
	}

	world = w
}
