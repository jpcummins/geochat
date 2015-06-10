package chat

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"time"
)

type DbConnection interface {
	Get() redis.Conn
	Publish(event *Event, zone *Zone) error
}

type RedisConnection struct {
	pool       *redis.Pool
	connection redis.Conn
}

func newRedisConnection(redisServer string) DbConnection {
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
	connection.connection = connection.Get()
	return connection
}

func (c *RedisConnection) Get() redis.Conn {
	return c.pool.Get()
}

func (c *RedisConnection) Publish(event *Event, zone *Zone) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err := c.connection.Do("PUBLISH", "zone_"+zone.id, eventJSON); err != nil {
		return err
	}
	return nil
}

type mockRedisConnection struct{}

func (mrc *mockRedisConnection) Close() error {
	println("Close()")
	return nil
}

func (mrc *mockRedisConnection) Err() error {
	println("Err()")
	return nil
}

func (mrc *mockRedisConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	fmt.Printf("Do(%s, ...)\n", commandName)
	return nil, nil
}

func (mrc *mockRedisConnection) Send(commandName string, args ...interface{}) error {
	fmt.Printf("Send(%s, ...)\n", commandName)
	return nil
}

func (mrc *mockRedisConnection) Flush() error {
	println("Flush()")
	return nil
}

func (mrc *mockRedisConnection) Receive() (reply interface{}, err error) {
	println("Receive()")
	return nil, nil
}

type mockConnection struct{}

func (mc *mockConnection) Get() redis.Conn {
	return &mockRedisConnection{}
}

func (mc *mockConnection) Publish(event *Event, zone *Zone) error {
	return event.Data.OnReceive(event, zone)
}
