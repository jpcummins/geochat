package chat

import (
	// "encoding/json"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/types"
	"strings"
	"sync"
)

// Zone represesnts a chat zone
type Zone struct {
	sync.RWMutex
	id       string
	boundary *ZoneBoundary
	geohash  string
	from     byte
	to       byte
	world    *World
	parent   *Zone
	left     *Zone
	right    *Zone
	maxUsers int
	isOpen   bool
	users    map[string]types.User
}

// ZoneBoundary provides the lat/long coordinates of the zone
type ZoneBoundary struct {
	SouthWestLat  float64 `json:"swlat"`
	SouthWestLong float64 `json:"swlong"`
	NorthEastLat  float64 `json:"nelat"`
	NorthEastLong float64 `json:"nelong"`
}

func newZone(world *World, geohash string, from byte, to byte, parent *Zone, maxUsers int) (*Zone, error) {
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
		users:    make(map[string]types.User),
		isOpen:   true,
		world:    world,
	}
	err := world.cache.SetZone(zone)
	return zone, err
}

func (z *Zone) createChildZones() {
	fromI := strings.Index(geohashmap, string(z.from))
	toI := strings.Index(geohashmap, string(z.to))

	if toI-fromI > 1 {
		split := (toI - fromI) / 2
		z.left, _ = newZone(z.world, z.geohash, z.from, geohashmap[fromI+split], z, z.maxUsers)
		z.right, _ = newZone(z.world, z.geohash, geohashmap[fromI+split+1], z.to, z, z.maxUsers)
	} else {
		z.left, _ = newZone(z.world, z.geohash+string(z.from), '0', 'z', z, z.maxUsers)
		z.right, _ = newZone(z.world, z.geohash+string(z.to), '0', 'z', z, z.maxUsers)
	}
}

func (z *Zone) ID() string {
	return z.id
}

func (z *Zone) Count() int {
	z.RLock()
	count := len(z.users)
	z.RUnlock()
	return count
}

func (z *Zone) IsOpen() bool {
	return z.Count() < z.maxUsers
}

func (z *Zone) AddUser(user types.User) error {
	z.Lock()
	z.users[user.ID()] = user
	err := z.world.cache.SetZone(z)
	z.Unlock()
	return err
}

func (z *Zone) RemoveUser(id string) error {
	z.Lock()
	delete(z.users, id)
	err := z.world.cache.SetZone(z)
	z.Unlock()
	return err
}

func (z *Zone) Broadcast(event types.Event) {
	z.RLock()
	for _, user := range z.users {
		user.Broadcast(event)
	}
	z.RUnlock()
}

// func (z *Zone) Users() map[string]*User {
// 	z.RLock()
// 	users := make(map[string]*User, len(z.users))
// 	for k, v := range z.users {
// 		users[k] = v
// 	}
// 	z.RUnlock()
// 	return users
// }

// func (z *Zone) split() {
// 	z.Lock()
// 	z.isOpen = false
// 	for _, user := range z.users {
// 		user.LeaveZone()
//
// 		zone, err := getOrCreateAvailableZone(user.lat, user.long)
// 		if err != nil {
// 			panic("Unable to create zone.")
// 		}
//
// 		user.JoinZone(zone)
// 	}
// 	z.Unlock()
// }
