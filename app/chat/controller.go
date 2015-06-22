package chat

import (
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

var App types.World

func Init(redisServer, worldID string) error {
	redisDB := db.NewRedisDB(redisServer)

	worlds := newWorlds(redisDB)
	world, err := worlds.World(worldID)
	if err != nil {
		return err
	}

	pubsub := db.NewRedisPubSub(worldID, redisDB)
	if world == nil {
		world, err = newWorld(worldID, redisDB, pubsub)
		if err != nil {
			return err
		}
		err = redisDB.SaveWorld(world.ServerJSON())
	}

	App = world
	return err
}
