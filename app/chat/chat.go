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
	"time"
	"os"
)

type Subscription struct {
	New  <-chan Event
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
		subscribe:   make(chan (chan<- Subscription), 10),
		unsubscribe: make(chan (<-chan Event), 10),
		publish:     make(chan Event, 10),
	}

	go zone.redisSubscribe()
	go zone.run()

	return zone
}

func (s Subscription) Unsubscribe() {
	s.Zone.unsubscribe <- s.New
}

func SubscribeToZone(geohash string) Subscription {
	zone, ok := zones[geohash]
	if !ok {
		zone = createZone(geohash)
		zones[geohash] = zone
	}

	subscription := make(chan Subscription)
	zone.subscribe <- subscription
	return <-subscription
}

func createPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
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

	if (redisServer == "") {
		redisServer = ":6379"
	}

	pool = createPool(redisServer)
	zones = make(map[string]*Zone)
}
