package chat

import (
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

var App types.World

func Init(redisServer, worldID string) error {
	redisDB := db.NewRedisDB(redisServer)

	world := &World{}
	found, err := redisDB.GetWorld(worldID, world)
	if err != nil {
		return err
	}

	if !found {
		pubsub := db.NewRedisPubSub(worldID, redisDB)
		world, err = newWorld(worldID, redisDB, pubsub)
		if err != nil {
			return err
		}
		err = redisDB.SetWorld(world)
	}

	App = world
	return err
}
