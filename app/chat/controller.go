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

	redisConnection := db.NewRedisConnection(redisServer)
	cache := cache.NewCache(redisConnection)

	factory := &Factory{}
	w, err := factory.NewWorld(worldID, cache)

	world = w
	return err
}
