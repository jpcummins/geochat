package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const rootWorldID string = "0"

type World struct {
	*sync.RWMutex
	*types.ServerWorldJSON
	root   types.Zone
	db     types.DB
	pubsub types.PubSub
	events types.Events
	users  types.Users
	zones  types.Zones
}

func newWorld(id string, db types.DB, ps types.PubSub, maxUsers int) (*World, error) {
	w := &World{
		ServerWorldJSON: &types.ServerWorldJSON{
			BaseServerJSON: &types.BaseServerJSON{
				ID: id,
			},
			MaxUsers: maxUsers,
		},
	}

	w.db = db
	w.pubsub = ps
	w.users = newUsers(id, db)
	w.zones = newZones(w, db)
	w.events = events.NewEvents(w)

	root, err := w.GetOrCreateZone(rootZoneID)
	if err != nil {
		return nil, err
	}

	w.root = root
	go w.manage() // It's a tough job.
	return w, nil
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
	return w.ServerWorldJSON.ID
}

func (w *World) MaxUsers() int {
	return w.ServerWorldJSON.MaxUsers
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
		zone, err = newZone(id, w, w.MaxUsers())
		if err != nil {
			return nil, err
		}

		if err := w.zones.Save(zone); err != nil {
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
	user := newUser(id, name, newLatLng(lat, lng), w)
	if err := w.Users().Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (w *World) NewServerEvent(data types.ServerEventData) types.ServerEvent {
	return w.events.NewServerEvent(generateEventID(), data)
}

func (w *World) NewClientEvent(data types.ClientEventData) types.ClientEvent {
	return w.events.NewClientEvent(generateEventID(), data)
}

func (w *World) Publish(event types.ServerEvent) error {
	if err := event.Data().BeforePublish(event); err != nil {
		return err
	}

	return w.pubsub.Publish(event)
}

func (w *World) ServerJSON() types.ServerJSON {
	return w.ServerWorldJSON
}

func (w *World) ClientJSON() types.ClientJSON {
	return nil
}

func (w *World) Update(json types.ServerJSON) error {
	worldJSON, ok := json.(*types.ServerWorldJSON)
	if !ok {
		return errors.New("Invalid json type.")
	}

	w.Lock()
	defer w.Unlock()
	w.ServerWorldJSON = worldJSON
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateEventID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10) + randomSequence(4)
}
