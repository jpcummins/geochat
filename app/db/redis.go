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

func (r *RedisDB) User(id string, worldID string) (*types.ServerUserJSON, error) {
	json := &types.ServerUserJSON{}
	err := r.getObject(getUserKey(id, worldID), json)
	return json, err
}

func (r *RedisDB) SaveUser(json types.ServerJSON) error {
	return r.setObject(getUserKey(json.Key(), json.WorldKey()), json)
}

func (r *RedisDB) Zone(id string, worldID string) (*types.ServerZoneJSON, error) {
	json := &types.ServerZoneJSON{}
	err := r.getObject(getZoneKey(id, worldID), json)
	return json, err
}

func (r *RedisDB) SaveZone(json types.ServerJSON) error {
	return r.setObject(getZoneKey(json.Key(), json.WorldKey()), json)
}

func (r *RedisDB) World(id string) (*types.ServerWorldJSON, error) {
	json := &types.ServerWorldJSON{}
	err := r.getObject(getWorldKey(id), json)
	return json, err
}

func (r *RedisDB) SaveWorld(json types.ServerJSON) error {
	return r.setObject(getWorldKey(json.Key()), json)
}

func (r *RedisDB) getObject(id string, v interface{}) error {
	data, err := redis.Bytes(r.connection.Do("GET", id))
	if err == redis.ErrNil {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
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

func getUserKey(id string, worldID string) string {
	return userPrefix + id + "_" + getWorldKey(worldID)
}

const zonePrefix = "zone_"

func getZoneKey(zoneID string, worldID string) string {
	return zonePrefix + zoneID + "_" + getWorldKey(worldID)
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

func (r *RedisPubSub) Publish(event types.ServerEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = r.db.connection.Do("PUBLISH", getWorldKey(r.worldID), string(bytes))
	return err
}

func (r *RedisPubSub) Subscribe() <-chan types.ServerEvent {
	ch := make(chan types.ServerEvent)
	go r.subscribe(ch)
	return ch
}

func (r *RedisPubSub) subscribe(ch chan types.ServerEvent) {
	psc := redis.PubSubConn{r.db.pool.Get()}
	defer psc.Close()
	psc.Subscribe(getWorldKey(r.worldID))

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var event types.ServerEvent
			if err := json.Unmarshal(v.Data, &event); err != nil {
				continue
			}
			ch <- event
		}
	}
}
