package chat

import (
	"github.com/jpcummins/geochat/app/db"
	// "github.com/jpcummins/geochat/app/events"
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
		world, err = newWorld(worldID, redisDB, events, 10)
		if err != nil {
			return nil, err
		}
		err = redisDB.SetWorld(world)
	}

	return world, err
}
