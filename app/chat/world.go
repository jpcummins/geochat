package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
	"strings"
	"sync"
)

type World struct {
	userMutex    sync.RWMutex
	zoneMutex    sync.RWMutex
	id           string
	root         types.Zone
	dependencies *Dependencies
	unsubscribe  chan bool
	users        map[string]types.User
	zones        map[string]types.Zone
}

const rootWorldID string = "0"

func newWorld(id string, dependencies *Dependencies) (*World, error) {
	world := &World{
		id:           id,
		dependencies: dependencies,
		unsubscribe:  make(chan bool),
		users:        make(map[string]types.User),
		zones:        make(map[string]types.Zone),
	}

	root, err := world.GetOrCreateZone(rootZoneID)
	if err != nil {
		return nil, err
	}

	world.root = root
	go world.manage() // It's a tough job.
	return world, nil
}

func (w *World) manage() {
	subscription := w.dependencies.PubSub().Subscribe()
	for {
		select {
		case event := <-subscription:
			event.SetWorld(w)
			event.Data().OnReceive(event)
		case <-w.unsubscribe:
			return
		}
	}
}

func (w *World) close() {
	if w != nil {
		w.unsubscribe <- true
		close(w.unsubscribe)
	}
}

func (w *World) ID() string {
	return w.id
}

func (w *World) Zone(id string) (types.Zone, error) {
	zone, found := w.localZone(id)
	if found {
		return zone, nil
	}
	return w.UpdateZone(id)
}

func (w *World) UpdateZone(id string) (types.Zone, error) {
	zone := &Zone{}
	found, err := w.dependencies.DB().GetZone(id, w, zone)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	w.localSetZone(zone)
	return zone, nil
}

func (w *World) SetZone(zone types.Zone) error {
	if err := w.dependencies.DB().SetZone(zone, w); err != nil {
		return err
	}

	w.localSetZone(zone)
	return nil
}

func (w *World) localZone(id string) (types.Zone, bool) {
	w.zoneMutex.RLock()
	defer w.zoneMutex.RUnlock()
	zone, found := w.zones[id]
	return zone, found
}

func (w *World) localSetZone(zone types.Zone) {
	w.zoneMutex.Lock()
	defer w.zoneMutex.Unlock()
	w.zones[zone.ID()] = zone
}

func (w *World) User(id string) (types.User, error) {
	user, found := w.localUser(id)
	if found {
		return user, nil
	}
	return w.UpdateUser(id)
}

func (w *World) UpdateUser(id string) (types.User, error) {
	user := &User{}
	found, err := w.dependencies.DB().GetUser(id, user)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	w.localSetUser(user)
	return user, nil
}

func (w *World) SetUser(user types.User) error {
	if err := w.dependencies.DB().SetUser(user); err != nil {
		return err
	}

	w.localSetUser(user)
	return nil
}

func (w *World) localUser(id string) (types.User, bool) {
	w.userMutex.RLock()
	defer w.userMutex.RUnlock()
	user, found := w.users[id]
	return user, found
}

func (w *World) localSetUser(user types.User) {
	w.userMutex.Lock()
	defer w.userMutex.Unlock()
	w.users[user.ID()] = user
}

func (w *World) GetOrCreateZone(id string) (types.Zone, error) {
	zone, err := w.Zone(id)
	if err != nil {
		return nil, err
	}

	if zone == nil {
		zone, err = newZone(id, w, 10)
		if err != nil {
			return nil, err
		}

		if err := w.SetZone(zone); err != nil {
			return nil, err
		}
	}

	return zone, nil
}

func (w *World) GetOrCreateZoneForUser(user types.User) (types.Zone, error) {
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

func (w *World) Publish(data types.EventData) error {
	event, eventErr := w.dependencies.Events().New("", data)
	if eventErr != nil {
		return eventErr
	}

	if publishErr := data.BeforePublish(event); publishErr != nil {
		return publishErr
	}

	return w.dependencies.PubSub().Publish(event)
}
