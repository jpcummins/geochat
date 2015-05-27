package chat

import (
	"os"
)

var LocalSubscribers *LocalSubscriptions
var connection *Connection
var world *World

func Init() {

	redisServer := os.Getenv("REDISTOGO_URL")

	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	connection = newConnection(redisServer)
	world = newWorld()
	LocalSubscribers = newSubscriptions()

	registerCommand(&command{
		name:    "addbot",
		usage:   "addbot (number of bots) (timeout in minutes) (geohash)",
		execute: addBot,
	})

	registerCommand(&command{
		name:    "flushdb",
		usage:   "",
		execute: resetRedis,
	})

	// registerCommand(&command{
	// 	name:    "join",
	// 	usage:   "",
	// 	execute: join,
	// })
}
