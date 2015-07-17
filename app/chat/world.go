package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/pubsub"
	"github.com/jpcummins/geochat/app/types"
	log "gopkg.in/inconshreveable/log15.v2"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const rootWorldID string = "0"

const defaultMaxUsers int = 30

const defaultMinUsers int = 10

const defaultSplitDelay time.Duration = time.Duration(30) * time.Second

type World struct {
	sync.RWMutex
	types.PubSubSerializable
	*types.WorldPubSubJSON

	root   types.Zone
	db     types.DB
	pubsub types.PubSub
	users  types.Users
	zones  types.Zones
	logger log.Logger
}

func newWorld(id string, db types.DB, ps types.PubSub, logger log.Logger) (*World, error) {
	w := &World{
		WorldPubSubJSON: &types.WorldPubSubJSON{
			ID:         id,
			MaxUsers:   defaultMaxUsers,
			MinUsers:   defaultMinUsers,
			SplitDelay: defaultSplitDelay,
		},
		db:     db,
		pubsub: ps,
		logger: logger.New("world", id),
	}

	w.users = newUsers(w, db)
	w.zones = newZones(w, db, logger)

	root, err := w.GetOrCreateZone(rootZoneID)
	if err != nil {
		w.logger.Crit("Error creating root zone", "root", rootZoneID, "error", err.Error())
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
	return w.WorldPubSubJSON.ID
}

func (w *World) MaxUsers() int {
	return w.WorldPubSubJSON.MaxUsers
}

func (w *World) MinUsers() int {
	return w.WorldPubSubJSON.MinUsers
}

func (w *World) SplitDelay() time.Duration {
	return w.WorldPubSubJSON.SplitDelay
}

func (w *World) DB() types.DB {
	return w.db
}

func (w *World) Users() types.Users {
	return w.users
}

func (w *World) Zone() types.Zone {
	return w.root
}

func (w *World) Zones() types.Zones {
	return w.zones
}

func (w *World) PubSub() types.PubSub {
	return w.pubsub
}

func (w *World) GetOrCreateZone(id string) (types.Zone, error) {
	zone, err := w.Zones().Zone(id)
	if err != nil {
		w.logger.Crit("Zone cache lookup error.", "zone", id, "error", err.Error())
		return nil, err
	}

	if zone == nil {
		zone, err = newZone(id, w, w.logger)
		if err != nil {
			w.logger.Crit("Error creating zone", "zone", id, "maxUsers", w.MaxUsers(), "error", err.Error())
			return nil, err
		}

		if err := w.zones.Save(zone); err != nil {
			w.logger.Crit("Error saving zone", "zone", id, "error", err.Error())
			return nil, err
		}
	}

	return zone, nil
}

func (w *World) FindOpenZone(root types.Zone, user types.User) (types.Zone, error) {
	for !root.IsOpen() {
		suffix := strings.TrimPrefix(user.Location().Geohash(), root.Geohash())

		if len(suffix) == 0 {
			w.logger.Crit("Unable to find an open zone", "user", user.ID())
			return nil, errors.New("Unable to find zone")
		}

		rightZone, err := w.GetOrCreateZone(root.RightZoneID())
		if err != nil {
			w.logger.Crit("Error retrieving right zone", "zone", root.RightZoneID())
			return nil, err
		}

		leftZone, err := w.GetOrCreateZone(root.LeftZoneID())
		if err != nil {
			w.logger.Crit("Error retrieving left zone", "zone", root.LeftZoneID())
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

func (w *World) Publish(data types.PubSubEventData) error {
	event := pubsub.NewEvent(generateEventID(), w, data)
	if err := event.Data().BeforePublish(event); err != nil {
		return err
	}
	return w.pubsub.Publish(event)
}

func (w *World) PubSubJSON() types.PubSubJSON {
	return w.WorldPubSubJSON
}

func (w *World) Update(json types.PubSubJSON) error {
	worldJSON, ok := json.(*types.WorldPubSubJSON)
	if !ok {
		return errors.New("Invalid json type.")
	}

	w.Lock()
	defer w.Unlock()
	w.WorldPubSubJSON = worldJSON
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
