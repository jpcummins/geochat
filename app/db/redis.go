package db

import (
	"encoding/json"
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

func (r *RedisDB) GetUser(id string, user types.User) (bool, error) {
	return r.getObject(getUserKey(id), user)
}

func (r *RedisDB) SetUser(user types.User) error {
	return r.setObject(getUserKey(user.ID()), user)
}

func (r *RedisDB) GetZone(id string, world types.World, zone types.Zone) (bool, error) {
	return r.getObject(getZoneKey(id, world.ID()), zone)
}

func (r *RedisDB) SetZone(zone types.Zone, world types.World) error {
	return r.setObject(getZoneKey(zone.ID(), world.ID()), zone)
}

func (r *RedisDB) GetWorld(id string, world types.World) (bool, error) {
	return r.getObject(getWorldKey(id), world)
}

func (r *RedisDB) SetWorld(world types.World) error {
	return r.setObject(getWorldKey(world.ID()), world)
}

func (r *RedisDB) getObject(id string, v interface{}) (bool, error) {
	data, err := redis.Bytes(r.connection.Do("GET", id))
	if err != nil {
		return data != nil, err
	}
	if data == nil {
		return false, nil
	}
	return true, json.Unmarshal(data, v)
}

func (r *RedisDB) setObject(id string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = r.connection.Do("SET", id, string(bytes))
	return err
}

const userPrefix = "user_"

func getUserKey(id string) string {
	return userPrefix + id
}

const zonePrefix = "zone_"

func getZoneKey(zoneID string, worldID string) string {
	return zonePrefix + zoneID + ":" + getWorldKey(worldID)
}

const worldPrefix = "world_"

func getWorldKey(id string) string {
	return worldPrefix + id
}

type RedisPubSub struct {
	world types.World
	db    *RedisDB
}

func NewRedisPubSub(world types.World, db *RedisDB) *RedisPubSub {
	return &RedisPubSub{world, db}
}

func (r *RedisPubSub) Publish(event types.Event) error {
	return nil
}

func (r *RedisPubSub) Subscribe() <-chan types.Event {
	return nil
}
