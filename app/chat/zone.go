package chat

import (
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

type Zone struct {
	sync.RWMutex
	*types.ServerZoneJSON
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
		ServerZoneJSON: &types.ServerZoneJSON{
			BaseServerJSON: &types.BaseServerJSON{
				ID:      id,
				WorldID: world.ID(),
			},
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

func (z *Zone) ID() string {
	return z.ServerZoneJSON.BaseServerJSON.ID
}

func (z *Zone) World() types.World {
	return z.world
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
	return z.ServerZoneJSON.MaxUsers
}

func (z *Zone) Count() int {
	z.RLock()
	defer z.RUnlock()
	return len(z.users)
}

func (z *Zone) IsOpen() bool {
	return z.ServerZoneJSON.IsOpen
}

func (z *Zone) SetIsOpen(isOpen bool) {
	z.ServerZoneJSON.IsOpen = isOpen
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

func (z *Zone) Broadcast(event types.ClientEvent) {
	z.RLock()
	defer z.RUnlock()

	for _, user := range z.users {
		user.Broadcast(event)
	}
}

func (z *Zone) ClientJSON() types.ClientJSON {
	return nil
}

func (z *Zone) ServerJSON() types.ServerJSON {
	z.RLock()
	defer z.RUnlock()
	z.ServerZoneJSON.UserIDs = make([]string, 0, len(z.users))
	for id := range z.users {
		z.ServerZoneJSON.UserIDs = append(z.ServerZoneJSON.UserIDs, id)
	}
	sort.Strings(z.ServerZoneJSON.UserIDs)
	return z.ServerZoneJSON
}

func (z *Zone) Update(js types.ServerJSON) error {
	json, ok := js.(*types.ServerZoneJSON)
	if !ok {
		return errors.New("Invalid json type.")
	}

	z.Lock()
	defer z.Unlock()
	z.ServerZoneJSON = json
	return nil
}

func (z *Zone) Join(user types.User) (types.ClientEvent, error) {
	joinEventData, err := events.NewJoin(z, user)
	if err != nil {
		return nil, err
	}
	event := z.world.NewServerEvent(joinEventData)
	return nil, z.world.Publish(event)
}

func (z *Zone) Message(user types.User, message string) (types.ClientEvent, error) {
	messageEventData, err := events.NewMessage(user, message)
	if err != nil {
		return nil, err
	}
	event := z.world.NewServerEvent(messageEventData)
	return nil, z.world.Publish(event)
}
