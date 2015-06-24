package pubsub

import (
	"errors"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const joinType types.PubSubEventType = "join"

type join struct {
	UserJSON *types.UserPubSubJSON `json:"user"`
	ZoneJSON *types.ZonePubSubJSON `json:"zone"`
	zone     types.Zone
	user     types.User
}

func Join(zone types.Zone, user types.User) (*join, error) {
	userJSON, ok := user.PubSubJSON().(*types.UserPubSubJSON)
	if !ok {
		return nil, errors.New("Unable to serialize UserPubSubJSON")
	}

	zoneJSON, ok := zone.PubSubJSON().(*types.ZonePubSubJSON)
	if !ok {
		return nil, errors.New("Unable to serialize ZonePubSubJSON")
	}

	return &join{
		UserJSON: userJSON,
		ZoneJSON: zoneJSON,
		zone:     zone,
		user:     user,
	}, nil
}

func (j *join) Type() types.PubSubEventType {
	return joinType
}

func (j *join) BeforePublish(e types.PubSubEvent) error {
	if j.user.Zone() != nil && j.user.Zone() != j.zone {
		// create and publish leave event
	}

	j.zone.AddUser(j.user)
	j.user.SetZone(j.zone)

	if err := e.World().Users().Save(j.user); err != nil {
		return err
	}

	return e.World().Zones().Save(j.zone)
}

func (j *join) OnReceive(e types.PubSubEvent) error {
	zone := e.World().Zones().FromCache(j.ZoneJSON.ID)
	if zone == nil {
		return nil
	}

	user := e.World().Users().FromCache(j.UserJSON.ID)
	if user == nil {
		var err error
		user, err = e.World().Users().FromDB(j.UserJSON.ID)
		if err != nil {
			return err
		}
	} else {
		user.Update(j.UserJSON)
	}
	zone.Update(j.ZoneJSON)
	zone.Broadcast(broadcast.Join(user, zone))
	return nil
}
