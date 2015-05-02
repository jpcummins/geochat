package chat

import (
	"github.com/revel/revel"
	"os"
)

var connection *Connection
var subscribers *Subscribers
var world *World

type Chat struct {
	*revel.Controller
	Connection  *Connection
	Subscribers *Subscribers
	World       *World
}

func Init() {

	redisServer := os.Getenv("REDISTOGO_URL")

	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	connection = newConnection(redisServer)
	world = newWorld()
	subscribers = newSubscribers()

	// registerCommand(&command{
	// 	name:    "addbot",
	// 	usage:   "addbot (number of bots) (timeout in minutes) (geohash)",
	// 	execute: addBot,
	// })

	// registerCommand(&command{
	// 	name:    "flushdb",
	// 	usage:   "",
	// 	execute: resetRedis,
	// })
}

func (c *Chat) Begin() revel.Result {
	c.Connection = connection
	c.Subscribers = subscribers
	c.World = world
	return nil
}

func init() {
	revel.InterceptMethod((*Chat).Begin, revel.BEFORE)
}
