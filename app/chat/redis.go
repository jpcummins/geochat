package chat

import (
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"time"
)

type RedisConnection struct {
	pool *redis.Pool
}

func newRedisConnection(redisServer string) *RedisConnection {
	connection := &RedisConnection{}
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redisurl.ConnectToURL(redisServer)
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
	connection.pool = pool
	return connection
}

func (c *RedisConnection) Get() redis.Conn {
	return c.pool.Get()
}
