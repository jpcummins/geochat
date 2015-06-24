package db

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jpcummins/geochat/app/pubsub"
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

func (r *RedisDB) User(id string, worldID string) (*types.UserPubSubJSON, error) {
	json := &types.UserPubSubJSON{}
	found, err := r.getObject(getUserKey(id, worldID), json)

	if !found {
		return nil, err
	}
	return json, err
}

func (r *RedisDB) SaveUser(json *types.UserPubSubJSON, worldID string) error {
	return r.setObject(getUserKey(json.ID, worldID), json)
}

func (r *RedisDB) Zone(id string, worldID string) (*types.ZonePubSubJSON, error) {
	json := &types.ZonePubSubJSON{}
	found, err := r.getObject(getZoneKey(id, worldID), json)

	if !found {
		return nil, err
	}

	return json, err
}

func (r *RedisDB) SaveZone(json *types.ZonePubSubJSON, worldID string) error {
	return r.setObject(getZoneKey(json.ID, worldID), json)
}

func (r *RedisDB) World(id string) (*types.WorldPubSubJSON, error) {
	json := &types.WorldPubSubJSON{}
	found, err := r.getObject(getWorldKey(id), json)

	if !found {
		return nil, err
	}

	return json, err
}

func (r *RedisDB) SaveWorld(json *types.WorldPubSubJSON) error {
	return r.setObject(getWorldKey(json.ID), json)
}

func (r *RedisDB) getObject(id string, v interface{}) (bool, error) {
	data, err := redis.Bytes(r.connection.Do("GET", id))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
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

func (r *RedisPubSub) Publish(event types.PubSubEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = r.db.connection.Do("PUBLISH", getWorldKey(r.worldID), string(bytes))
	return err
}

func (r *RedisPubSub) Subscribe() <-chan types.PubSubEvent {
	ch := make(chan types.PubSubEvent)
	go r.subscribe(ch)
	return ch
}

func (r *RedisPubSub) subscribe(ch chan types.PubSubEvent) {
	psc := redis.PubSubConn{r.db.pool.Get()}
	defer psc.Close()
	psc.Subscribe(getWorldKey(r.worldID))

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var event pubsub.Event
			if err := json.Unmarshal(v.Data, &event); err != nil {
				println("Error unmarshaling from pubsub: ", err.Error())
				continue
			}
			ch <- &event
		}
	}
}
