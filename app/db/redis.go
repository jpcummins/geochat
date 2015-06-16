package db

import (
	// "encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jpcummins/geochat/app/types"
	"github.com/soveran/redisurl"
	"time"
)

type Redis struct {
	pool       *redis.Pool
	connection redis.Conn
}

func NewRedisConnection(redisServer string) *Redis {
	connection := &Redis{}
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

func (r *Redis) Publish(event types.Event) error {
	// eventJSON, err := json.Marshal(event)
	// if err != nil {
	// 	return err
	// }
	// if _, err := r.connection.Do("PUBLISH", "zone_"+event.Zone().ID(), eventJSON); err != nil {
	// 	return err
	// }
	return nil
}

func (r *Redis) Subscribe(world types.World) <-chan types.Event {
	return make(<-chan types.Event)
}

func (r *Redis) GetUser(id string) (types.User, error) {
	return nil, nil
}

func (r *Redis) SetUser(types.User) error {
	return nil
}

func (r *Redis) GetZone(id string) (types.Zone, error) {
	return nil, nil
}

func (r *Redis) SetZone(types.Zone) error {
	return nil
}

func (r *Redis) GetWorld(id string) (types.World, error) {
	return nil, nil
}

func (r *Redis) SetWorld(types.World) error {
	return nil
}

// func (u *UserCache) redisGet(id string) (types.User, bool) {
// 	c := Redis.Get()
// 	defer c.Close()
//
// 	usersJSON, err := redis.String(c.Do("HGET", "users", id))
// 	if err != nil {
// 		return nil, false
// 	}
//
// 	var user UserCache
// 	if err := json.Unmarshal([]byte(usersJSON), &user); err != nil {
// 		return nil, false
// 	}
// 	return &user, true
// }

// func (u *UserCache) redisSet(user *User) {
// 	c := Redis.Get()
// 	defer c.Close()
// 	eventJSON, err := json.Marshal(user)
// 	if err != nil {
// 		return
// 	}
// 	c.Do("HSET", "users", user.id, eventJSON)
// }

//
// type mockRedisConnection struct{}
//
// func (mrc *mockRedisConnection) Close() error {
// 	println("Close()")
// 	return nil
// }
//
// func (mrc *mockRedisConnection) Err() error {
// 	println("Err()")
// 	return nil
// }
//
// func (mrc *mockRedisConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
// 	fmt.Printf("Do(%s, ...)\n", commandName)
// 	return nil, nil
// }
//
// func (mrc *mockRedisConnection) Send(commandName string, args ...interface{}) error {
// 	fmt.Printf("Send(%s, ...)\n", commandName)
// 	return nil
// }
//
// func (mrc *mockRedisConnection) Flush() error {
// 	println("Flush()")
// 	return nil
// }
//
// func (mrc *mockRedisConnection) Receive() (reply interface{}, err error) {
// 	println("Receive()")
// 	return nil, nil
// }
//
// type mockConnection struct{}
//
// func (mc *mockConnection) Get() redis.Conn {
// 	return &mockRedisConnection{}
// }
//
// func (mc *mockConnection) Publish(event *Event, zone *Zone) error {
// 	return event.Data.OnReceive(event, zone)
// }
