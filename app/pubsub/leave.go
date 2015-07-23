package pubsub

import (
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
	l.user.SetZoneID("")

	if err := e.World().Users().Save(l.user); err != nil {
		return err
	}

	return e.World().Zones().Save(l.zone)
}

func (l *leave) OnReceive(e types.PubSubEvent) error {

	// refresh cache. At some point this should be optimize. I don't think a
	// db hit is nessesary.
	zone, err := e.World().Zones().FromDB(l.ZoneID)
	if err != nil {
		return err
	}

	user, err := e.World().Users().FromDB(l.UserID)
	if err != nil {
		return err
	}

	zone.Broadcast(broadcast.Leave(user))
	return nil
}
