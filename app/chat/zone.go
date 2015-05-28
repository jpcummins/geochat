package chat

import (
	"encoding/json"
	"fmt"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/garyburd/redigo/redis"
	"strings"
)

// Zone represesnts a chat zone
type Zone struct {
	id          string
	boundary    *ZoneBoundary
	geohash     string
	from        byte
	to          byte
	parent      *Zone
	left        *Zone
	right       *Zone
	count       int
	maxUsers    int
	publish     chan *Event
	subscribers []*Subscription
}

// ZoneBoundary provides the lat/long coordinates of the zone
type ZoneBoundary struct {
	SouthWestLat  float64 `json:"swlat"`
	SouthWestLong float64 `json:"swlong"`
	NorthEastLat  float64 `json:"nelat"`
	NorthEastLong float64 `json:"nelong"`
}

// ZoneJSON is passed to the client when a websocket connection is established
type ZoneJSON struct {
	ID          string          `json:"id"`
	Boundary    *ZoneBoundary   `json:"boundary"`
	Archive     *Archive        `json:"archive"`
	Subscribers []*Subscription `json:"subscribers"`
}

// Type implements EventType, which is used to provide Event.UnmarshalJSON a
// hint on how to unmarshal Zone JSON.
func (z *Zone) Type() string {
	return "zone"
}

// OnReceive implements EventType. This method is called when a "zone" event is
// received from Redis.
func (z *Zone) OnReceive(event *Event) error {
	return nil
}

// MarshalJSON overrides Go's default JSON marshaling method so that this object
// can be marshaled/unmarshaled as the type `zoneJSON`
func (z *Zone) MarshalJSON() ([]byte, error) {
	js := &ZoneJSON{
		ID:          z.GetID(),
		Boundary:    z.GetBoundary(),
		Archive:     z.GetArchive(50),
		Subscribers: z.GetSubscribers(),
	}
	return json.Marshal(js)
}

func newZone(geohash string, from byte, to byte, parent *Zone, maxUsers int) *Zone {
	sw := gh.Decode(geohash + string(from))
	ne := gh.Decode(geohash + string(to))

	zone := &Zone{
		id:      geohash + ":" + string(from) + string(to),
		geohash: geohash,
		boundary: &ZoneBoundary{
			SouthWestLat:  sw.SouthWest().Lat(),
			SouthWestLong: sw.SouthWest().Lng(),
			NorthEastLat:  ne.NorthEast().Lat(),
			NorthEastLong: ne.NorthEast().Lng(),
		},
		from:     from,
		to:       to,
		parent:   parent,
		maxUsers: maxUsers,
		publish:  make(chan *Event, 10),
	}

	go zone.redisSubscribe() // subscribe to zone's redis channel
	go zone.redisPublish()   // publishes publish events to redis channel

	return zone
}

func (z *Zone) createChildZones() {
	fromI := strings.Index(geohashmap, string(z.from))
	toI := strings.Index(geohashmap, string(z.to))

	if toI-fromI > 1 {
		split := (toI - fromI) / 2
		z.left = newZone(z.geohash, z.from, geohashmap[fromI+split], z, z.maxUsers)
		z.right = newZone(z.geohash, geohashmap[fromI+split+1], z.to, z, z.maxUsers)
	} else {
		z.left = newZone(z.geohash+string(z.from), '0', 'z', z, z.maxUsers)
		z.right = newZone(z.geohash+string(z.to), '0', 'z', z, z.maxUsers)
	}
}

func (z *Zone) redisSubscribe() {
	psc := redis.PubSubConn{Conn: connection.Get()}
	defer psc.Close()
	psc.Subscribe("zone_" + z.id)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var event Event
			if err := json.Unmarshal(v.Data, &event); err != nil {
				println("Error unmarshaling event: ", err.Error())
			}

			if err := event.Data.OnReceive(&event); err != nil {
				println("Error executing event: ", err.Error())
			}
		}
	}
}

func (z *Zone) redisPublish() {
	c := connection.Get()
	defer c.Close()
	for {
		select {
		case event := <-z.publish:
			eventJSON, _ := json.Marshal(event)
			c.Do("LPUSH", "zone_"+z.id, eventJSON)
			c.Do("PUBLISH", "zone_"+z.id, eventJSON)
		}
	}
}

// GetID returns the zone identifier
func (z *Zone) GetID() string {
	return z.id
}

// GetSubscribers returns all of the zone's subscribers
func (z *Zone) GetSubscribers() []*Subscription {
	return z.subscribers
}

// GetArchive returns the last N events in the zone
func (z *Zone) GetArchive(maxEvents int) *Archive {
	c := connection.Get()
	defer c.Close()

	archiveJSON, err := redis.Strings(c.Do("LRANGE", "zone_"+z.id, 0, maxEvents-1))
	if err != nil {
		println("Unable to get archive:", err.Error())
		return nil
	}

	return newArchive(archiveJSON)
}

// GetBoundary returns the zone's coordinates
func (z *Zone) GetBoundary() *ZoneBoundary {
	return z.boundary
}

// Publish publishes an event to the zone
func (z *Zone) Publish(event *Event) {
	z.publish <- event
}

func (z *Zone) join(s *Subscription) {
	z.subscribers = append(z.subscribers, s)
	fmt.Printf("Subscribers: %+v\n", z.subscribers)
	incrementZoneSubscriptionCounts(z) // bubble up the count
}

func (z *Zone) leave(s *Subscription) {
	for i, x := range z.subscribers {
		if x.id == s.id {
			copy(z.subscribers[i:], z.subscribers[i+1:])
			z.subscribers[len(z.subscribers)-1] = nil
			z.subscribers = z.subscribers[:len(z.subscribers)-1]
			decrementZoneSubscriptionCounts(z) // bubble up the count
			return
		}
	}
}
