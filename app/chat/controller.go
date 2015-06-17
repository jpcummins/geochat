package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

var factory types.Factory

func Init(redisServer, worldID string) (types.World, error) {
	redisConnection := db.NewRedisDB(redisServer)
	cache := cache.NewCache(redisConnection)
	pubsub, err := db.NewRedisPubSub(worldID, redisConnection)

	if err != nil {
		return nil, err
	}

	factory := &Factory{}
	return factory.NewWorld(worldID, cache, pubsub)
}
