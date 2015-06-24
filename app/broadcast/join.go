package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const joinType types.BroadcastEventType = "join"

type join struct {
	UserJSON interface{} `json:"user"`
	ZoneJSON interface{} `json:"zone,omitempty"`
	user     types.User
	zone     types.Zone
}

func Join(user types.User, zone types.Zone) *join {
	return &join{
		UserJSON: user.BroadcastJSON(),
		user:     user,
		zone:     zone,
	}
}

func (e *join) Type() types.BroadcastEventType {
	return joinType
}

func (e *join) BeforeBroadcastToUser(user types.User, event types.BroadcastEvent) error {
	if e.user == user {
		e.ZoneJSON = e.zone.BroadcastJSON()
	}
	return nil
}
