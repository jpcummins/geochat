package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const joinType types.BroadcastEventType = "join"

type join struct {
	UserJSON *types.UserBroadcastJSON `json:"user"`
	ZoneJSON *types.ZoneBroadcastJSON `json:"zone"`
}

func Join(user types.User, zone types.Zone) *join {
	return &join{
		UserJSON: user.BroadcastJSON().(*types.UserBroadcastJSON),
		ZoneJSON: zone.BroadcastJSON().(*types.ZoneBroadcastJSON),
	}
}

func (e *join) Type() types.BroadcastEventType {
	return joinType
}
