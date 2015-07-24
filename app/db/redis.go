package db

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jpcummins/geochat/app/pubsub"
	"github.com/jpcummins/geochat/app/types"
	"github.com/soveran/redisurl"
	log "gopkg.in/inconshreveable/log15.v2"
	"strings"
	"sync"
	"time"
)

type RedisDB struct {
	sync.RWMutex
	pool   *redis.Pool
	logger log.Logger
}

func NewRedisDB(redisServer string, logger log.Logger) *RedisDB {
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
	connection.logger = logger.New("redis", redisServer)
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

	if err != nil {
		r.logger.Error("Error retrieving zone", "id", id, "world", worldID, "error", err.Error())
		return nil, err
	}

	if !found {
		r.logger.Info("Zone not found", "id", id, "world", worldID)
		return nil, nil
	}

	return json, err
}

func (r *RedisDB) SaveZone(json *types.ZonePubSubJSON, worldID string) error {
	json.LastModified = time.Now()
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
	r.Lock()
	defer r.Unlock()

	connection := r.pool.Get()

	data, err := redis.Bytes(connection.Do("GET", id))
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
	r.Lock()
	defer r.Unlock()

	connection := r.pool.Get()
	_, err = connection.Do("SET", id, string(bytes))
	return err
}

func (r *RedisDB) SaveUsersAndZones(users []*types.UserPubSubJSON, zones []*types.ZonePubSubJSON, worldID string) error {
	r.Lock()
	defer r.Unlock()
	connection := r.pool.Get()

	// Gather the object keys
	keys := make([]string, len(users)+len(zones))
	for _, user := range users {
		keys = append(keys, getUserKey(user.ID, worldID))
	}
	for _, zone := range zones {
		keys = append(keys, getZoneKey(zone.ID, worldID))
	}

	// Watch the keys for changes
	connection.Send("WATCH", strings.Join(keys, " "))
	connection.Send("MULTI")

	for _, user := range users {
		bytes, err := json.Marshal(user)
		if err != nil {
			r.logger.Error("Unable to marshal user", "user", user.ID)
			return err
		}
		connection.Send("SET", getUserKey(user.ID, worldID), bytes)
	}

	for _, zone := range zones {
		bytes, err := json.Marshal(zone)
		if err != nil {
			r.logger.Error("Unable to marshal zone", "zone", zone.ID)
			return err
		}
		connection.Send("SET", getZoneKey(zone.ID, worldID), bytes)
	}

	if _, err := connection.Do("EXEC"); err != nil {
		r.logger.Error("Error executing save", "error", err.Error())
		return err
	}

	_, err := connection.Do("UNWATCH")

	if err != nil {
		r.logger.Error("Error unwatching", "error", err.Error())
	}

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
	return &RedisPubSub{
		worldID: worldID,
		db:      db,
	}
}

func (pubsub *RedisPubSub) Publish(event types.PubSubEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	pubsub.db.Lock()
	defer pubsub.db.Unlock()

	connection := pubsub.db.pool.Get()
	_, err = connection.Do("PUBLISH", getWorldKey(pubsub.worldID), string(bytes))
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
