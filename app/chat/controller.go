package chat

import (
	"github.com/jpcummins/geochat/app/db"
	"github.com/jpcummins/geochat/app/types"
	log "gopkg.in/inconshreveable/log15.v2"
)

var App types.World

func Init(redisServer, worldID string) error {
	logger := log.Root()
	redisDB := db.NewRedisDB(redisServer, logger)
	pubsub := db.NewRedisPubSub(worldID, redisDB)

	worlds := newWorlds(redisDB, pubsub, logger)
	world, err := worlds.World(worldID)
	if err != nil {
		logger.Crit("Error when looking up world in cache.", "error", err.Error())
		return err
	}

	if world == nil {
		world, err = newWorld(worldID, redisDB, pubsub, logger)
		if err != nil {
			logger.Crit("Error when creating new world.", "error", err.Error())
			return err
		}
		if err := worlds.Save(world); err != nil {
			logger.Crit("Error saving world.", "error", err.Error())
			return err
		}
	}

	App = world
	return err
}
