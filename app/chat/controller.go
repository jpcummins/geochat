package chat

import (
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
)

func Init(redisServer, worldID string) (types.World, error) {
	redisDB := db.NewRedisDB(redisServer)

	world := &World{}
	found, err := redisDB.GetWorld(worldID, world)
	if err != nil {
		return nil, err
	}

	if !found {
		pubsub := db.NewRedisPubSub(worldID, redisDB)
		world, err = newWorld(worldID, redisDB, pubsub)
		if err != nil {
			return nil, err
		}
		err = redisDB.SetWorld(world)
	}

	return world, err
}
