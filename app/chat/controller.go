package chat

import (
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
		err = redisDB.SaveWorld(world.ServerJSON())
	}

	App = world
	return err
}
