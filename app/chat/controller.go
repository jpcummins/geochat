package chat

import (
	"os"
)

var UserCache *Users
var connection *RedisConnection
var world *World

func Init() {

	redisServer := os.Getenv("REDISTOGO_URL")

	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	UserCache = NewUsers()
	connection = newRedisConnection(redisServer)
	world = newWorld()
	world.root.initialize()

	registerCommand(&command{
		name:    "addbot",
		usage:   "addbot (number of bots) (timeout in minutes)",
		execute: addBot,
	})

	registerCommand(&command{
		name:    "addbot2",
		usage:   "addbot2 (number of bots) (timeout in minutes)",
		execute: addBot2,
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
