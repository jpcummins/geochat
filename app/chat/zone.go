package chat

import (
	"encoding/json"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type Zone struct {
	Zonehash string        `json:"zonehash"`
	Boundary *ZoneBoundary `json:"boundary"`
	geohash  string        `json:"-"`
	from     byte          `json:"-"`
	to       byte          `json:"-"`
	parent   *Zone         `json:"-"`
	left     *Zone         `json:"-"`
	right    *Zone         `json:"-"`
	count    int           `json:"-"`
	maxUsers int           `json:"-"`
	publish  chan *Event   `json:"-"`
}

type ZoneBoundary struct {
	SouthWestLat  float64
	SouthWestLong float64
	NorthEastLat  float64
	NorthEastLong float64
}

func (z *Zone) Type() string {
	return "zone"
}

func newZone(geohash string, from byte, to byte, parent *Zone, maxUsers int) *Zone {
	zone := &Zone{
		Zonehash: geohash + ":" + string(from) + string(to),
		Boundary: GetBoundary(geohash, from, to),
		geohash:  geohash,
		from:     from,
		to:       to,
		parent:   parent,
		maxUsers: maxUsers,
		publish:  make(chan *Event, 10),
	}
	return zone
}

func (z *Zone) setCount(count int) {
	if z.count == 0 && count > 0 {
		// The following goroutines terminate on their own when count == 0
		go z.redisSubscribe() // subscribe to redis channel and publishes events
		go z.redisPublish()   // publishes publish events to redis channel
	}

	if count >= 0 {
		z.count = count
	}
}

func GetBoundary(geohash string, from, to byte) *ZoneBoundary {
	sw := gh.Decode(geohash + string(from))
	ne := gh.Decode(geohash + string(to))
	return &ZoneBoundary{
		SouthWestLat:  sw.SouthWest().Lat(),
		SouthWestLong: sw.SouthWest().Lng(),
		NorthEastLat:  ne.NorthEast().Lat(),
		NorthEastLong: ne.NorthEast().Lng(),
	}
}

func (z *Zone) GetSubscribers() []*Subscription {
	ch := make(chan interface{})
	subscribers.getSubscriptionsForZone <- ch
	ch <- z
	return (<-ch).([]*Subscription)
}

func (z *Zone) createChildZones() {
	from_i := strings.Index(geohashmap, string(z.from))
	to_i := strings.Index(geohashmap, string(z.to))

	if to_i-from_i > 1 {
		split := (to_i - from_i) / 2
		z.left = newZone(z.geohash, z.from, geohashmap[from_i+split], z, z.maxUsers)
		z.right = newZone(z.geohash, geohashmap[from_i+split+1], z.to, z, z.maxUsers)
	} else {
		z.left = newZone(z.geohash+string(z.from), '0', 'z', z, z.maxUsers)
		z.right = newZone(z.geohash+string(z.to), '0', 'z', z, z.maxUsers)
	}
}

func (z *Zone) GetArchive(maxEvents int) *Archive {
	c := connection.Get()
	defer c.Close()

	archiveJson, err := redis.Strings(c.Do("LRANGE", "zone_"+z.Zonehash, 0, maxEvents-1))
	if err != nil {
		println("unable to get archive:", err.Error())
		return nil
	}

	return newArchive(archiveJson)
}

func (z *Zone) Publish(event *Event) {
	z.publish <- event
}

func (z *Zone) redisSubscribe() {
	psc := redis.PubSubConn{connection.Get()}
	defer psc.Close()
	psc.Subscribe("zone_" + z.Zonehash)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:

			if z.count == 0 {
				return
			}

			var event Event
			if err := json.Unmarshal(v.Data, &event); err != nil {
				continue
			}
			subscribers.PublishEventToZone(&event, z)
		}
	}
}

func (z *Zone) redisPublish() {
	c := connection.Get() // Not sure if a long lived redis connection is a good idea
	defer c.Close()
	for {

		if z.count == 0 {
			return
		}

		select {
		case event := <-z.publish:
			eventJson, _ := json.Marshal(event)

			c.Do("LPUSH", "zone_"+z.Zonehash, eventJson)
			c.Do("PUBLISH", "zone_"+z.Zonehash, eventJson)
		}
	}
}
