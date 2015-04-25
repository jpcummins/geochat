package chat

import (
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"os"
	"time"
)

func createPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redisurl.ConnectToURL(server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	pool          *redis.Pool
	subscriptions map[string]*Subscription
	maxRoomSize   int32
	world         *Zone
)

func init() {

	redisServer := os.Getenv("REDISTOGO_URL")

	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	maxRoomSize = 1
	pool = createPool(redisServer)
	subscriptions = make(map[string]*Subscription)

	world = createZone("", '0', 'z', nil)

	registerCommand(&command{
		name:    "addbot",
		usage:   "addbot (number of bots) (timeout in minutes) (geohash)",
		execute: addBot,
	})

	registerCommand(&command{
		name:    "fuckeverything",
		usage:   "",
		execute: resetRedis,
	})
}
