package chat

import (
	"encoding/json"
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/types"
	"sort"
	"strings"
	"sync"
)

type zoneJSON struct {
	ID       string   `json:"id"`
	UserIDs  []string `json:"user_ids"`
	IsOpen   bool     `json:"is_open"`
	MaxUsers int      `json:"max_users"`
}

type Zone struct {
	sync.RWMutex
	*zoneJSON
	southWest types.LatLng
	northEast types.LatLng
	geohash   string
	from      byte
	to        byte
	parent    *Zone
	left      *Zone
	right     *Zone
	users     map[string]types.User
}

func newZone(world types.World, id string) (*Zone, error) {
	geohash, from, to, err := validateZoneID(id)
	if err != nil {
		return nil, err
	}

	zone := &Zone{
		zoneJSON: &zoneJSON{
			ID:       geohash + ":" + string(from) + string(to),
			IsOpen:   true,
			MaxUsers: world.MaxUsersForNewZones(),
		},
		southWest: gh.Decode(geohash + string(from)).SouthWest(),
		northEast: gh.Decode(geohash + string(to)).NorthEast(),
		geohash:   geohash,
		from:      from,
		to:        to,
		parent:    nil, // parent,
		users:     make(map[string]types.User),
	}
	return zone, nil
}

func validateZoneID(id string) (geohash string, from byte, to byte, err error) {
	split := strings.Split(id, ":")

	if len(split) != 2 || len(split[1]) != 2 {
		err = errors.New("Invalid id")
	}

	// TODO: Additional validation needed
	geohash = split[0]
	from = split[1][0]
	to = split[1][1]
	return
}

func (z *Zone) MarshalJSON() ([]byte, error) {
	z.RLock()
	z.zoneJSON.UserIDs = make([]string, 0, len(z.users))
	for id, _ := range z.users {
		z.zoneJSON.UserIDs = append(z.zoneJSON.UserIDs, id)
	}
	z.RUnlock()
	sort.Strings(z.zoneJSON.UserIDs)
	return json.Marshal(z.zoneJSON)
}

func (z *Zone) ID() string {
	return z.zoneJSON.ID
}

func (z *Zone) SouthWest() types.LatLng {
	return z.southWest
}

func (z *Zone) NorthEast() types.LatLng {
	return z.northEast
}

func (z *Zone) Geohash() string {
	return z.geohash
}

func (z *Zone) From() byte {
	return z.from
}

func (z *Zone) To() byte {
	return z.to
}

func (z *Zone) Parent() types.Zone {
	return z.parent
}

func (z *Zone) Left() types.Zone {
	return z.left
}

func (z *Zone) Right() types.Zone {
	return z.right
}

func (z *Zone) MaxUsers() int {
	return z.zoneJSON.MaxUsers
}

func (z *Zone) Count() int {
	z.RLock()
	count := len(z.users)
	z.RUnlock()
	return count
}

func (z *Zone) IsOpen() bool {
	return z.zoneJSON.IsOpen
}

func (z *Zone) SetIsOpen(isOpen bool) {
	z.zoneJSON.IsOpen = isOpen
}

func (z *Zone) AddUser(user types.User) {
	z.Lock()
	z.users[user.ID()] = user
	z.Unlock()
}

func (z *Zone) RemoveUser(id string) {
	z.Lock()
	delete(z.users, id)
	z.Unlock()
}

func (z *Zone) Broadcast(event types.Event) {
	z.RLock()
	for _, user := range z.users {
		user.Broadcast(event)
	}
	z.RUnlock()
}

// func (z *Zone) createChildZones() {
// 	fromI := strings.Index(geohashmap, string(z.from))
// 	toI := strings.Index(geohashmap, string(z.to))
//
// 	if toI-fromI > 1 {
// 		split := (toI - fromI) / 2
// 		z.left = newZone(z.geohash, z.from, geohashmap[fromI+split], z, z.maxUsers)
// 		z.right = newZone(z.geohash, geohashmap[fromI+split+1], z.to, z, z.maxUsers)
// 	} else {
// 		z.left = newZone(z.geohash+string(z.from), '0', 'z', z, z.maxUsers)
// 		z.right = newZone(z.geohash+string(z.to), '0', 'z', z, z.maxUsers)
// 	}
// }
