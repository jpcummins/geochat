package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
)

var chat types.Chat

func Init(redisServer, worldID string) (types.Chat, error) {
	redisConnection := db.NewRedisDB(redisServer)
	cache := cache.NewCache(redisConnection)

	world, err := cache.World(worldID)
	if err != nil {
		return nil, err
	}

	pubsub, err := db.NewRedisPubSub(worldID, redisConnection)
	if err != nil {
		return nil, err
	}

	events := events.NewEventFactory(world)
	return newChat(cache, pubsub, events, 2)
}
