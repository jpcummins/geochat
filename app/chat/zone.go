package chat

import (
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/pubsub"
	"github.com/jpcummins/geochat/app/types"
	"strings"
	"sync"
)

const rootZoneID = ":0z"

const geohashmap = "0123456789bcdefghjkmnpqrstuvwxyz"

type Zone struct {
	sync.RWMutex
	types.PubSubSerializable
	types.BroadcastSerializable
	*types.ZonePubSubJSON

	world        types.World
	southWest    types.LatLng
	northEast    types.LatLng
	geohash      string
	from         string
	to           string
	parentZoneID string
	leftZoneID   string
	rightZoneID  string
}

func newZone(id string, world types.World, maxUsers int) (*Zone, error) {
	geohash, from, to, err := validateZoneID(id)
	if err != nil {
		return nil, err
	}

	southWest := gh.Decode(geohash + from).SouthWest()
	northEast := gh.Decode(geohash + to).NorthEast()

	zone := &Zone{
		ZonePubSubJSON: &types.ZonePubSubJSON{
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
	return z.ZonePubSubJSON.ID
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
	return z.ZonePubSubJSON.MaxUsers
}

func (z *Zone) Count() int {
	z.RLock()
	defer z.RUnlock()
	return len(z.ZonePubSubJSON.UserIDs)
}

func (z *Zone) IsOpen() bool {
	return z.ZonePubSubJSON.IsOpen
}

func (z *Zone) SetIsOpen(isOpen bool) {
	z.ZonePubSubJSON.IsOpen = isOpen
}

func (z *Zone) AddUser(user types.User) {
	z.Lock()
	defer z.Unlock()
	z.ZonePubSubJSON.UserIDs = append(z.ZonePubSubJSON.UserIDs, user.ID())
}

func (z *Zone) RemoveUser(id string) {
	z.Lock()
	defer z.Unlock()

	users := z.ZonePubSubJSON.UserIDs
	for i := range users {
		if users[i] == id {
			z.ZonePubSubJSON.UserIDs = append(users[:i], users[i+1:]...)
			return
		}
	}
}

func (z *Zone) Broadcast(eventData types.BroadcastEventData) {
	z.RLock()
	defer z.RUnlock()

	for _, id := range z.ZonePubSubJSON.UserIDs {
		if user, err := z.World().Users().User(id); user != nil && err == nil {
			user.Broadcast(eventData)
		}
	}
}

func (z *Zone) BroadcastJSON() interface{} {
	z.RLock()
	defer z.RUnlock()
	json := &types.ZoneBroadcastJSON{
		ID:        z.ID(),
		Users:     make(map[string]*types.UserBroadcastJSON),
		SouthWest: z.southWest.BroadcastJSON().(*types.LatLngJSON),
		NorthEast: z.northEast.BroadcastJSON().(*types.LatLngJSON),
	}
	for _, id := range z.ZonePubSubJSON.UserIDs {
		if user, err := z.World().Users().User(id); err == nil {
			json.Users[id] = user.BroadcastJSON().(*types.UserBroadcastJSON)
		}
	}
	return json
}

func (z *Zone) PubSubJSON() types.PubSubJSON {
	return z.ZonePubSubJSON
}

func (z *Zone) Update(js types.PubSubJSON) error {
	json, ok := js.(*types.ZonePubSubJSON)
	if !ok {
		return errors.New("Unable to serialize to ZonePubSubJSON.")
	}

	z.Lock()
	defer z.Unlock()
	z.ZonePubSubJSON = json
	return nil
}

func (z *Zone) Join(user types.User) (types.BroadcastEvent, error) {
	if user.Zone() != nil && user.Zone() != z {
		if _, err := user.Zone().Leave(user); err != nil {
			return nil, err
		}
	}

	data, err := pubsub.Join(z, user)
	if err != nil {
		return nil, err
	}
	return nil, z.world.Publish(data)
}

func (z *Zone) Leave(user types.User) (types.BroadcastEvent, error) {
	data, err := pubsub.Leave(user, z)
	if err != nil {
		return nil, err
	}
	return nil, z.world.Publish(data)
}

func (z *Zone) Message(user types.User, message string) (types.BroadcastEvent, error) {
	data, err := pubsub.Message(user, z, message)
	if err != nil {
		return nil, err
	}
	return nil, z.world.Publish(data)
}
