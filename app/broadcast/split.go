package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const splitType types.BroadcastEventType = "split"

type split struct {
	PreviousZone interface{} `json:"previous_zone"`
	NextZone     interface{} `json:"next_zone"`
	zones        map[string]types.Zone
}

func Split(previousZoneID string, zones map[string]types.Zone) *split {
	return &split{
		PreviousZone: zones[previousZoneID].BroadcastJSON(),
		zones:        zones,
	}
}

func (e *split) Type() types.BroadcastEventType {
	return splitType
}

func (e *split) BeforeBroadcastToUser(user types.User, event types.BroadcastEvent) (bool, error) {
	e.NextZone = e.zones[user.Zone().ID()].BroadcastJSON()
	return true, nil
}
