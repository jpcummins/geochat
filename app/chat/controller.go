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

	pubsub := db.NewRedisPubSub(worldID, redisDB)
	if found {
		if err := world.init(redisDB, pubsub); err != nil {
			return err
		}
	} else {
		world = newWorld(worldID)
		if err := world.init(redisDB, pubsub); err != nil {
			return err
		}
		err = redisDB.SetWorld(world)
	}

	App = world
	return err
}
