package chat

/*

TODO:
* hash the geohash for security?
* Add geolocation detection in the browser
* Perhaps I should remove the user filter to make things more explicit

*/

import (
	"container/list"
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"os"
	"time"
)

type Subscription struct {
	Events chan *Event
    User *User
	Zone *Zone
}

var zones map[string]*Zone

func FindZone(zone string) (z *Zone, ok bool) {
	z, ok = zones[zone]
	return
}

func createZone(geohash string) *Zone {
	zone := &Zone{
		Geohash:     geohash,
		subscribers: list.New(),
		publish:     make(chan *Event, 10),
	}

	go zone.redisSubscribe()
	go zone.run()
	return zone
}

func (s *Subscription) Unsubscribe() {
    for subscriber := s.Zone.subscribers.Front(); subscriber != nil; subscriber = subscriber.Next() {
        s.Zone.subscribers.Remove(subscriber)
        s.Zone.Publish(NewEvent(&Leave{User: s.User}))
    }
}

func SubscribeToZone(geohash string, user *User) (*Subscription, *Zone) {
	zone, ok := zones[geohash]
	if !ok {
		zone = createZone(geohash)
		zones[geohash] = zone
	}

    subscription := &Subscription{
        User: user,
        Zone: zone,
        Events: make(chan *Event, 10),
    }

	zone.subscribers.PushBack(subscription)
    zone.Publish(NewEvent(&Join{User: user}))
	return subscription, zone
}

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
	pool *redis.Pool
)

func init() {

	redisServer := os.Getenv("REDISTOGO_URL")

	if redisServer == "" {
		redisServer = "redis://localhost:6379"
	}

	pool = createPool(redisServer)
	zones = make(map[string]*Zone)
}
