package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

var chat types.Chat

func Init(redisServer, worldID string) (types.Chat, error) {
	redisConnection := db.NewRedisDB(redisServer)
	cache := cache.NewCache(redisConnection)
	pubsub, err := db.NewRedisPubSub(worldID, redisConnection)

	if err != nil {
		return nil, err
	}

	return newChat(cache, pubsub, 2)
}
