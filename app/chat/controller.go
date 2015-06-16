package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

var world types.World

func Init(redisServer, worldID string) error {
	// redisServer := os.Getenv("REDISTOGO_URL")
	// if redisServer == "" {
	// 	redisServer = "redis://localhost:6379"
	// }

	redisConnection := db.NewRedisDB(redisServer)
	cache := cache.NewCache(redisConnection)
	pubsub, err := db.NewRedisPubSub(worldID, redisConnection)

	if err != nil {
		return err
	}

	factory := &Factory{}
	w, err := factory.NewWorld(worldID, cache, pubsub)

	world = w
	return err
}
