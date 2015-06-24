package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

var App types.World

func Init(redisServer, worldID string, maxUsers int) error {
	redisDB := db.NewRedisDB(redisServer)
	pubsub := db.NewRedisPubSub(worldID, redisDB)

	worlds := newWorlds(redisDB, pubsub)
	world, err := worlds.World(worldID)
	if err != nil {
		return err
	}

	if world == nil {
		world, err = newWorld(worldID, redisDB, pubsub, maxUsers)
		if err != nil {
			return err
		}

		worldJSON, ok := world.PubSubJSON().(*types.WorldPubSubJSON)
		if !ok {
			return errors.New("Unable to serialize WorldPubSubJSON")
		}
		err = redisDB.SaveWorld(worldJSON)
	}

	App = world
	return err
}
