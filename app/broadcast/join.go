package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const joinType types.BroadcastEventType = "join"

type join struct {
	UserJSON *types.UserBroadcastJSON `json:"user"`
	user     types.User
	zone     types.Zone
}

func Join(user types.User, zone types.Zone) *join {
	return &join{
		UserJSON: user.BroadcastJSON().(*types.UserBroadcastJSON),
		user:     user,
		zone:     zone,
	}
}

func (e *join) Type() types.BroadcastEventType {
	return joinType
}

func (e *join) BeforeBroadcastToUser(user types.User, event types.BroadcastEvent) (bool, error) {
	if e.user == user {
		user.Broadcast(Zone(e.zone))
		return false, nil
	}
	return true, nil
}
