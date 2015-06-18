package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
	"strings"
)

type World struct {
	id              string
	root            types.Zone
	chat            types.Chat
	maxUsersPerZone int
	subscribe       <-chan types.Event
	unsubscribe     chan bool
}

func newWorld(id string, chat types.Chat, maxUsersPerZone int) (*World, error) {
	world := &World{
		id:              id,
		chat:            chat,
		maxUsersPerZone: maxUsersPerZone,
		subscribe:       chat.PubSub().Subscribe(),
		unsubscribe:     make(chan bool),
	}

	root, err := world.GetOrCreateZone(":0z")
	if err != nil {
		return nil, err
	}

	world.root = root
	go world.manage() // It's a tough job.
	return world, nil
}

func (w *World) manage() {
	for {
		select {
		case event := <-w.subscribe:
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
	return w.chat.Cache().Zone(id)
}

func (w *World) SetZone(zone types.Zone) error {
	return w.chat.Cache().SetZone(zone)
}

func (w *World) UpdateZone(id string) (types.Zone, error) {
	return w.chat.Cache().UpdateZone(id)
}

func (w *World) User(id string) (types.User, error) {
	return w.chat.Cache().User(id)
}

func (w *World) SetUser(user types.User) error {
	return w.chat.Cache().SetUser(user)
}

func (w *World) UpdateUser(id string) (types.User, error) {
	return w.chat.Cache().UpdateUser(id)
}

func (w *World) GetOrCreateZone(id string) (types.Zone, error) {
	zone, err := w.Zone(id)
	if err != nil {
		return nil, err
	}

	if zone == nil {
		zone, err = newZone(id, w, w.maxUsersPerZone)
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
	event, eventErr := w.chat.Events().New("", data)
	if eventErr != nil {
		return eventErr
	}

	if publishErr := data.BeforePublish(event); publishErr != nil {
		return publishErr
	}

	return w.chat.PubSub().Publish(event)
}
