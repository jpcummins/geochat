package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type World struct {
	id     string
	root   types.Zone
	db     types.DB
	pubsub types.PubSub
	events types.Events
	users  types.Users
	zones  types.Zones
}

const rootWorldID string = "0"

func newWorld(id string, db types.DB, ps types.PubSub) (*World, error) {
	world := &World{
		id:     id,
		db:     db,
		pubsub: ps,
	}

	world.users = newUsers(db, world)
	world.zones = newZones(db, world)
	world.events = events.NewEvents(world)

	root, err := world.GetOrCreateZone(rootZoneID)
	if err != nil {
		return nil, err
	}

	world.root = root
	go world.manage() // It's a tough job.
	return world, nil
}

func (w *World) manage() {
	subscription := w.pubsub.Subscribe()
	for {
		select {
		case event := <-subscription:
			event.SetWorld(w)
			event.Data().OnReceive(event)
		}
	}
}

func (w *World) ID() string {
	return w.id
}

func (w *World) Users() types.Users {
	return w.users
}

func (w *World) Zones() types.Zones {
	return w.zones
}

func (w *World) GetOrCreateZone(id string) (types.Zone, error) {
	zone, err := w.Zones().Zone(id)
	if err != nil {
		return nil, err
	}

	if zone == nil {
		zone, err = newZone(id, w, 10)
		if err != nil {
			return nil, err
		}

		if err := w.zones.SetZone(zone); err != nil {
			return nil, err
		}
	}

	return zone, nil
}

func (w *World) FindOpenZone(user types.User) (types.Zone, error) {
	root := w.root
	for !root.IsOpen() {
		suffix := strings.TrimPrefix(user.Location().Geohash(), root.Geohash())

		if len(suffix) == 0 {
			return nil, errors.New("Unable to find zone")
		}

		rightZone, err := w.GetOrCreateZone(root.RightZoneID())
		if err != nil {
			return nil, err
		}

		leftZone, err := w.GetOrCreateZone(root.LeftZoneID())
		if err != nil {
			return nil, err
		}

		if rightZone.Geohash() == root.Geohash() {
			if suffix[0] < rightZone.From()[0] {
				root = leftZone
			} else {
				root = rightZone
			}
		} else {
			if suffix[0] < rightZone.Geohash()[len(rightZone.Geohash())-1] {
				root = leftZone
			} else {
				root = rightZone
			}
		}
	}

	return root, nil
}

func (w *World) NewUser(id string, name string, lat float64, lng float64) (types.User, error) {
	user := newUser(id, name, newLatLng(lat, lng))
	if err := w.Users().SetUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (w *World) Publish(data types.EventData) error {
	id := strconv.FormatInt(time.Now().UnixNano(), 10) + randomSequence(4)

	event, eventErr := w.events.New(id, data)
	if eventErr != nil {
		return eventErr
	}

	if publishErr := data.BeforePublish(event); publishErr != nil {
		return publishErr
	}

	return w.pubsub.Publish(event)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
