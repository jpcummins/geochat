package chat

import (
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/pubsub"
	"github.com/jpcummins/geochat/app/types"
	log "gopkg.in/inconshreveable/log15.v2"
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
	logger       log.Logger
}

func newZone(id string, world types.World, logger log.Logger) (*Zone, error) {
	geohash, from, to, err := validateZoneID(id)
	if err != nil {
		return nil, err
	}

	southWest := gh.Decode(geohash + from).SouthWest()
	northEast := gh.Decode(geohash + to).NorthEast()

	zone := &Zone{
		ZonePubSubJSON: &types.ZonePubSubJSON{
			ID:     id,
			IsOpen: true,
		},
		world:     world,
		southWest: newLatLng(southWest.Lat(), southWest.Lng()),
		northEast: newLatLng(northEast.Lat(), northEast.Lng()),
		geohash:   geohash,
		from:      from,
		to:        to,
		logger:    logger.New("zone", id),
	}

	// Calculate left, right, and parent IDs

	// so gross
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

	// oh god, what am I doing?
	if from == "0" && to == "z" {
		if geohash == "" {
			zone.parentZoneID = ""
		} else {
			parentFromOrTo := geohash[len(geohash)-1:]
			parentI := strings.Index(geohashmap, parentFromOrTo)
			if parentI%2 == 0 {
				zone.parentZoneID = geohash[:len(geohash)-1] + ":" + parentFromOrTo + string(geohashmap[parentI+1])
			} else {
				zone.parentZoneID = geohash[:len(geohash)-1] + ":" + string(geohashmap[parentI-1]) + parentFromOrTo
			}
		}
	} else {
		distance := (toI - fromI + 1) * 2
		if fromI%distance == 0 {
			zone.parentZoneID = geohash + ":" + from + string(geohashmap[fromI+distance-1])
		} else {
			zone.parentZoneID = geohash + ":" + string(geohashmap[toI-distance+1]) + to
		}
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

func (z *Zone) Count() int {
	z.RLock()
	defer z.RUnlock()
	return len(z.ZonePubSubJSON.UserIDs)
}

func (z *Zone) UserIDs() []string {
	ids := make([]string, len(z.ZonePubSubJSON.UserIDs))
	copy(ids, z.ZonePubSubJSON.UserIDs)
	return ids
}

func (z *Zone) IsOpen() bool {
	return z.ZonePubSubJSON.IsOpen
}

func (z *Zone) SetIsOpen(isOpen bool) {
	z.ZonePubSubJSON.IsOpen = isOpen
	z.World().Zones().UpdateCache(z)
}

func (z *Zone) AddUser(user types.User) {
	z.Lock()
	defer z.Unlock()

	// If the user is already here, don't add.
	users := z.ZonePubSubJSON.UserIDs
	for i := range users {
		if users[i] == user.ID() {
			return
		}
	}

	z.ZonePubSubJSON.UserIDs = append(z.ZonePubSubJSON.UserIDs, user.ID())
	z.World().Zones().UpdateCache(z)
}

func (z *Zone) RemoveUser(id string) {
	z.Lock()
	defer z.Unlock()

	users := z.ZonePubSubJSON.UserIDs
	for i := range users {
		if users[i] == id {
			z.ZonePubSubJSON.UserIDs = append(users[:i], users[i+1:]...)
			z.World().Zones().UpdateCache(z)
			return
		}
	}
}

func (z *Zone) Broadcast(eventData types.BroadcastEventData) error {
	z.RLock()
	defer z.RUnlock()

	for _, id := range z.ZonePubSubJSON.UserIDs {
		if user, err := z.World().Users().User(id); user != nil && err == nil {
			user.Broadcast(eventData)
		}
	}

	return nil
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
	z.World().Zones().UpdateCache(z)
	return nil
}

func (z *Zone) Join(user types.User) error {
	if user.Zone() != nil && user.Zone() != z {
		if err := user.Zone().Leave(user); err != nil {
			return err
		}
	}

	data, err := pubsub.Join(z, user)
	if err != nil {
		return err
	}
	return z.world.Publish(data)
}

func (z *Zone) Leave(user types.User) error {
	data, err := pubsub.Leave(user, z)
	if err != nil {
		return err
	}
	return z.world.Publish(data)
}

func (z *Zone) Message(user types.User, message string) error {
	data, err := pubsub.Message(user, z, message)
	if err != nil {
		return err
	}
	return z.world.Publish(data)
}

func (z *Zone) Split() (map[string]types.Zone, error) {
	// Close the zone and save
	z.SetIsOpen(false)
	z.World().Zones().Save(z)

	// Update the user and zone objects
	users := make(map[string]types.User)
	zones := make(map[string]types.Zone)

	for _, userID := range z.ZonePubSubJSON.UserIDs {
		user, err := z.World().Users().User(userID)
		if err != nil {
			z.logger.Crit("User cache lookup error", "error", err.Error(), "user", userID)
			return nil, err
		}

		newZone, err := z.world.FindOpenZone(z, user)
		if err != nil {
			z.logger.Crit("Unable to find an open zone", "currentZone", z.ID(), "user", user.ID())
			return nil, err
		}
		user.SetZone(newZone)
		newZone.AddUser(user)

		users[userID] = user
		zones[newZone.ID()] = newZone
	}

	// clear the zone subscriber list
	z.ZonePubSubJSON.UserIDs = z.ZonePubSubJSON.UserIDs[:0]

	// Bulk save new zones, current zone, and users.
	zones[z.ID()] = z
	usersSlice := make([]*types.UserPubSubJSON, 0, len(users))
	zonesSlice := make([]*types.ZonePubSubJSON, 0, len(zones))

	for _, user := range users {
		usersSlice = append(usersSlice, user.PubSubJSON().(*types.UserPubSubJSON))
	}

	for _, zone := range zones {
		zonesSlice = append(zonesSlice, zone.PubSubJSON().(*types.ZonePubSubJSON))
	}

	if err := z.world.DB().SaveUsersAndZones(usersSlice, zonesSlice, z.world.ID()); err != nil {
		z.logger.Crit("Error saving users and/or zones", "error", err.Error())
		return nil, err
	}
	return zones, nil
}

func (z *Zone) Merge() error {
	leftZone, err := z.world.GetOrCreateZone(z.LeftZoneID())

	if err != nil {
		z.logger.Error("Error retrieving left zone", "leftZone", z.LeftZoneID(), "error", err.Error())
		return err
	}

	if leftZone == nil {
		err := errors.New("Left zone (" + z.LeftZoneID() + ") is nil")
		z.logger.Error(err.Error())
		return err
	}

	rightZone, err := z.world.GetOrCreateZone(z.RightZoneID())

	if err != nil {
		z.logger.Error("Error retrieving right zone", "rightZone", z.RightZoneID(), "error", err.Error())
	}

	if rightZone == nil {
		err := errors.New("Right zone (" + z.RightZoneID() + ")")
		z.logger.Error(err.Error())
		return err
	}

	users := make([]*types.UserPubSubJSON, 0, z.Count()+leftZone.Count()+rightZone.Count())

	for _, userID := range leftZone.UserIDs() {
		user, err := z.world.Users().User(userID)
		if err != nil {
			z.logger.Error("Error retrieving user", "user", userID, "error", err.Error())
			return err
		}
		leftZone.RemoveUser(userID)
		z.AddUser(user)
		user.SetZone(z)
		users = append(users, user.PubSubJSON().(*types.UserPubSubJSON))
	}

	for _, userID := range rightZone.UserIDs() {
		user, err := z.world.Users().User(userID)
		if err != nil {
			z.logger.Error("Error retrieving user", "user", userID, "error", err.Error())
			return err
		}
		rightZone.RemoveUser(userID)
		z.AddUser(user)
		user.SetZone(z)
		users = append(users, user.PubSubJSON().(*types.UserPubSubJSON))
	}

	z.SetIsOpen(true)

	zones := make([]*types.ZonePubSubJSON, 3)
	zones[0] = z.PubSubJSON().(*types.ZonePubSubJSON)
	zones[1] = leftZone.PubSubJSON().(*types.ZonePubSubJSON)
	zones[2] = rightZone.PubSubJSON().(*types.ZonePubSubJSON)

	return z.world.DB().SaveUsersAndZones(users, zones, z.world.ID())
}
