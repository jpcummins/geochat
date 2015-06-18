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
	world, _ := cache.World(worldID)
	pubsub, err := db.NewRedisPubSub(worldID, redisConnection)
	events := events.NewEventFactory(world)

	if err != nil {
		return nil, err
	}

	return newChat(cache, pubsub, events, 2)
}
