package chat

import (
	"encoding/json"
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/garyburd/redigo/redis"
	"strings"
	"sync"
)

// Zone represesnts a chat zone
type Zone struct {
	sync.RWMutex
	id            string
	boundary      *ZoneBoundary
	geohash       string
	from          byte
	to            byte
	parent        *Zone
	left          *Zone
	right         *Zone
	count         int
	maxUsers      int
	isOpen        bool
	publish       chan *Event
	archive       chan *Event
	announceJoin  chan *User
	announceLeave chan *User
	users         map[string]*User
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
	ID       string           `json:"id"`
	Boundary *ZoneBoundary    `json:"boundary"`
	Archive  *Archive         `json:"archive"`
	Users    map[string]*User `json:"users"`
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
		ID:       z.GetID(),
		Boundary: z.GetBoundary(),
		Archive:  z.GetArchive(50),
		Users:    z.GetUsers(),
	}
	json, err := json.Marshal(js)
	return json, err
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
		users:    make(map[string]*User),
		isOpen:   true,
	}
	return zone
}

func (z *Zone) isInitialized() bool {
	return z.publish != nil
}

func (z *Zone) initialize() {
	z.publish = make(chan *Event, 10)
	z.archive = make(chan *Event, 10)
	z.announceJoin = make(chan *User, 10)
	z.announceLeave = make(chan *User, 10)

	c := connection.Get()
	defer c.Close()

	ids, err := redis.Strings(c.Do("SMEMBERS", "users_"+z.id))
	if err != nil {
		panic(errors.New("Unable to download users for zone " + z.id))
	}

	for _, id := range ids {
		user, found := UserCache.Get(id)
		if found {
			IncrementZoneSubscriptionCounts(z) // optimize this at some point.
			z.addUser(user)
		}
	}

	go z.redisSubscribe() // subscribe to zone's redis channel
	go z.redisPublish()   // publishes publish events to redis channel
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
			c.Do("PUBLISH", "zone_"+z.id, eventJSON)
		case event := <-z.archive:
			eventJSON, _ := json.Marshal(event)
			c.Do("LPUSH", "archive_"+z.id, eventJSON)
		case user := <-z.announceJoin:
			c.Do("SADD", "users_"+z.id, user.GetID())
		case user := <-z.announceLeave:
			c.Do("SREM", "users_"+z.id, user.GetID())
		}
	}
}

// GetID returns the zone identifier
func (z *Zone) GetID() string {
	if z == nil {
		return ""
	}
	return z.id
}

// GetArchive returns the last N events in the zone
func (z *Zone) GetArchive(maxEvents int) *Archive {
	c := connection.Get()
	defer c.Close()

	archiveJSON, err := redis.Strings(c.Do("LRANGE", "archive_"+z.id, 0, maxEvents-1))
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

func (z *Zone) GetUsers() map[string]*User {
	z.RLock()
	users := make(map[string]*User, len(z.users))
	for k, v := range z.users {
		users[k] = v
	}
	z.RUnlock()
	return users
}

func (z *Zone) join(u *User) {
	IncrementZoneSubscriptionCounts(z)
	z.announceJoin <- u
	z.Publish(NewEvent(&Join{User: u}))

	if z.count > z.maxUsers {
		z.split()
	}
}

func (z *Zone) leave(u *User) {
	DecrementZoneSubscriptionCounts(z)
	z.announceLeave <- u
	z.Publish(NewEvent(&Leave{UserID: u.GetID(), ZoneID: u.zone.id}))
}

func (z *Zone) addUser(u *User) {
	z.Lock()
	z.users[u.GetID()] = u
	z.Unlock()
}

func (z *Zone) delUser(userID string) {
	z.Lock()
	delete(z.users, userID)
	z.Unlock()
}

// Publish publishes an event to the zone's Redis channel
func (z *Zone) Publish(event *Event) {
	z.publish <- event
}

// Broadcast sends an event to all local users in the zone
func (z *Zone) broadcastEvent(event *Event) {
	z.RLock()
	for _, user := range z.users {
		for _, connection := range user.connections {
			connection.Events <- event
		}
	}
	z.RUnlock()
}

func (z *Zone) archiveEvent(event *Event) {
	z.archive <- event
}

func (z *Zone) split() {
	z.Lock()
	z.isOpen = false
	for _, user := range z.users {
		user.LeaveZone()

		zone, err := getOrCreateAvailableZone(user.lat, user.long)
		if err != nil {
			panic("Unable to create zone.")
		}

		user.JoinZone(zone)
	}
	z.Unlock()
}
