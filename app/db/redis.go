package db

import (
	// "encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jpcummins/geochat/app/types"
	"github.com/soveran/redisurl"
	"time"
)

type RedisDB struct {
	pool       *redis.Pool
	connection redis.Conn
}

func NewRedisDB(redisServer string) *RedisDB {
	connection := &RedisDB{}
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
	connection.connection = pool.Get()
	return connection
}

func (r *RedisDB) GetUser(id string) (types.User, error) {
	return nil, nil
}

func (r *RedisDB) SetUser(types.User) error {
	return nil
}

func (r *RedisDB) GetZone(id string) (types.Zone, error) {
	return nil, nil
}

func (r *RedisDB) SetZone(types.Zone) error {
	return nil
}

func (r *RedisDB) GetWorld(id string) (types.World, error) {
	return nil, nil
}

func (r *RedisDB) SetWorld(types.World) error {
	return nil
}

type RedisPubSub struct {
	worldID string
	db      *RedisDB
}

func NewRedisPubSub(worldID string, db *RedisDB) (*RedisPubSub, error) {
	ps := &RedisPubSub{worldID, db}
	return ps, nil
}

func (r *RedisPubSub) Publish(event types.Event) error {
	return nil
}

func (r *RedisPubSub) Subscribe() <-chan types.Event {
	return nil
}
