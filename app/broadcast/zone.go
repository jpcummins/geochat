package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const zoneType types.BroadcastEventType = "zone"

type zone struct {
	ZoneJSON *types.ZoneBroadcastJSON `json:"zone"`
}

func Zone(z types.Zone) *zone {
	return &zone{
		ZoneJSON: z.BroadcastJSON().(*types.ZoneBroadcastJSON),
	}
}

func (e *zone) Type() types.BroadcastEventType {
	return zoneType
}

func (e *zone) BeforeBroadcastToUser(user types.User, event types.BroadcastEvent) (bool, error) {
	return true, nil
}
