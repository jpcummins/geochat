package chat

import (
	"encoding/json"
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
	"sort"
	"strings"
	"sync"
)

const rootZoneID = ":0z"

const geohashmap = "0123456789bcdefghjkmnpqrstuvwxyz"

type zoneJSON struct {
	ID       string   `json:"id"`
	UserIDs  []string `json:"user_ids"`
	IsOpen   bool     `json:"is_open"`
	MaxUsers int      `json:"max_users"`
}

type Zone struct {
	sync.RWMutex
	*zoneJSON
	world        types.World
	southWest    types.LatLng
	northEast    types.LatLng
	geohash      string
	from         string
	to           string
	parentZoneID string
	leftZoneID   string
	rightZoneID  string
	users        map[string]types.User
}

func newZone(id string, world types.World, maxUsers int) (*Zone, error) {
	geohash, from, to, err := validateZoneID(id)
	if err != nil {
		return nil, err
	}

	southWest := gh.Decode(geohash + from).SouthWest()
	northEast := gh.Decode(geohash + to).NorthEast()

	zone := &Zone{
		zoneJSON: &zoneJSON{
			ID:       id,
			IsOpen:   true,
			MaxUsers: maxUsers,
		},
		world:     world,
		southWest: newLatLng(southWest.Lat(), southWest.Lng()),
		northEast: newLatLng(northEast.Lat(), northEast.Lng()),
		geohash:   geohash,
		from:      from,
		to:        to,
		users:     make(map[string]types.User),
	}

	// Calculate left, right, and parent IDs
	fromI := strings.Index(geohashmap, from)
	toI := strings.Index(geohashmap, to)
	if toI-fromI > 1 {
		split := (toI - fromI) / 2
		zone.leftZoneID = geohash + ":" + from + string(geohashmap[fromI+split])
		zone.rightZoneID = geohash + ":" + string(geohashmap[fromI+split+1]) + to
	} else {
		zone.leftZoneID = geohash + from + rootZoneID
		zone.rightZoneID = geohash + to + rootZoneID
	}

	return zone, nil
}

func validateZoneID(id string) (geohash string, from string, to string, err error) {
	split := strings.Split(id, ":")

	if len(split) != 2 || len(split[1]) != 2 {
		return "", "", "", errors.New("Invalid id")
	}

	// TODO: Additional validation needed
	geohash = split[0]
	from = string(split[1][0])
	to = string(split[1][1])
	return
}

func (z *Zone) MarshalJSON() ([]byte, error) {
	z.RLock()
	defer z.RUnlock()
	z.zoneJSON.UserIDs = make([]string, 0, len(z.users))
	for id := range z.users {
		z.zoneJSON.UserIDs = append(z.zoneJSON.UserIDs, id)
	}
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

func (z *Zone) From() string {
	return z.from
}

func (z *Zone) To() string {
	return z.to
}

func (z *Zone) ParentZoneID() string {
	return z.parentZoneID
}

func (z *Zone) LeftZoneID() string {
	return z.leftZoneID
}

func (z *Zone) RightZoneID() string {
	return z.rightZoneID
}

func (z *Zone) MaxUsers() int {
	return z.zoneJSON.MaxUsers
}

func (z *Zone) Count() int {
	z.RLock()
	defer z.RUnlock()
	return len(z.users)
}

func (z *Zone) IsOpen() bool {
	return z.zoneJSON.IsOpen
}

func (z *Zone) SetIsOpen(isOpen bool) {
	z.zoneJSON.IsOpen = isOpen
}

func (z *Zone) AddUser(user types.User) {
	z.Lock()
	defer z.Unlock()
	z.users[user.ID()] = user
}

func (z *Zone) RemoveUser(id string) {
	z.Lock()
	defer z.Unlock()
	delete(z.users, id)
}

func (z *Zone) Broadcast(event types.Event) {
	z.RLock()
	defer z.RUnlock()

	for _, user := range z.users {
		user.Broadcast(event)
	}
}

func (z *Zone) Join(user types.User) error {
	joinEventData, err := events.NewJoin(z, user)
	if err != nil {
		return err
	}
	return z.world.Publish(joinEventData)
}

func (z *Zone) Message(user types.User, message string) error {
	messageEventData, err := events.NewMessage(user, z, message)
	if err != nil {
		return err
	}
	return z.world.Publish(messageEventData)
}
