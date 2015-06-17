package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
	"strings"
)

type World struct {
	id        string
	root      types.Zone
	chat      *Chat
	subscribe <-chan types.Event
}

func newWorld(id string, chat *Chat) (*World, error) {
	world := &World{
		id:        id,
		chat:      chat,
		subscribe: chat.pubsub.Subscribe(),
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
			event.Data().OnReceive(event)
		}
	}
}

func (w *World) ID() string {
	return w.id
}

func (w *World) Zone(id string) (types.Zone, error) {
	return w.chat.cache.Zone(id)
}

func (w *World) SetZone(zone types.Zone) error {
	return w.chat.cache.SetZone(zone)
}

func (w *World) SetUser(user types.User) error {
	return w.chat.cache.SetUser(user)
}

func (w *World) GetOrCreateZone(id string) (types.Zone, error) {
	zone, err := w.Zone(id)
	if err != nil {
		return nil, err
	}

	if zone == nil {
		zone, err = newZone(id, w.id, w.chat.maxUsersForNewZones)
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

func (w *World) Publish(event types.Event) error {
	if err := event.Data().BeforePublish(event); err != nil {
		return err
	}
	return w.chat.pubsub.Publish(event)
}
