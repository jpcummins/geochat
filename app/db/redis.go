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
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return data != nil, err
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
	worldID string
	db      *RedisDB
}

func NewRedisPubSub(worldID string, db *RedisDB) *RedisPubSub {
	return &RedisPubSub{worldID, db}
}

func (r *RedisPubSub) Publish(event types.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = r.db.connection.Do("PUBLISH", getWorldKey(r.worldID), string(bytes))
	return err
}

func (r *RedisPubSub) Subscribe() <-chan types.Event {
	ch := make(chan types.Event)
	go r.subscribe(ch)
	return ch
}

func (r *RedisPubSub) subscribe(ch chan types.Event) {
	psc := redis.PubSubConn{r.db.pool.Get()}
	defer psc.Close()
	psc.Subscribe(getWorldKey(r.worldID))

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var event types.Event
			if err := json.Unmarshal(v.Data, &event); err != nil {
				continue
			}
			ch <- event
		}
	}
}
