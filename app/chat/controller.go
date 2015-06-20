package chat

import (
	"github.com/jpcummins/geochat/app/db"
	// "github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
)

func Init(redisServer, worldID string) (types.World, error) {
	redisDB := db.NewRedisDB(redisServer)

	dependencies := &Dependencies{
		db: redisDB,
		// pubsub:              db.NewRedisPubSub(world, redisDB),
		// events:              events.NewEventFactory(world),
		maxUsersForNewZones: 10,
	}

	world := &World{}
	found, err := redisDB.GetWorld(worldID, world)
	if err != nil {
		return nil, err
	}

	if !found {
		return newWorld(worldID, dependencies)
	}

	return world, nil
}
