package pubsub

import (
	"errors"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const leaveType types.PubSubEventType = "leave"

type leave struct {
	UserID string `json:"user_id"`
	ZoneID string `json:"zone_id"`
	user   types.User
	zone   types.Zone
}

func Leave(user types.User, zone types.Zone) (*leave, error) {
	l := &leave{
		UserID: user.ID(),
		ZoneID: zone.ID(),
		user:   user,
		zone:   zone,
	}
	return l, nil
}

func (l *leave) Type() types.PubSubEventType {
	return leaveType
}

func (l *leave) BeforePublish(e types.PubSubEvent) error {
	l.zone.RemoveUser(l.UserID)
	l.user.SetZone(nil)

	if err := e.World().Users().Save(l.user); err != nil {
		return err
	}

	return e.World().Zones().Save(l.zone)
}

func (l *leave) OnReceive(e types.PubSubEvent) error {
	zone := e.World().Zones().FromCache(l.ZoneID)
	if zone == nil {
		return nil
	}

	user := e.World().Users().FromCache(l.UserID)
	if user == nil {
		return errors.New("Unable to find user: " + l.UserID)
	}

	user.SetZone(nil)
	zone.RemoveUser(l.UserID)
	zone.Broadcast(broadcast.Leave(user))
	return nil
}
