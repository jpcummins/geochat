package chat

import (
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
)

var chat types.Chat

func Init(redisServer, worldID string) (types.Chat, error) {
	redisDB := db.NewRedisDB(redisServer)

	world := &World{}
	_, err := redisDB.GetWorld(worldID, world)
	if err != nil {
		return nil, err
	}

	return &Chat{
		db:                  redisDB,
		pubsub:              db.NewRedisPubSub(world, redisDB),
		events:              events.NewEventFactory(world),
		world:               world,
		maxUsersForNewZones: 10,
	}, nil
}
