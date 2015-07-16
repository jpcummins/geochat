package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const splitType types.BroadcastEventType = "split"

type split struct {
	PreviousZone interface{} `json:"previous_zone"`
	Zone         interface{} `json:"zone"`
}

func Split(previousZone types.Zone, zone types.Zone) *split {
	return &split{
		PreviousZone: previousZone.BroadcastJSON(),
		Zone:         zone.BroadcastJSON(),
	}
}

func (e *split) Type() types.BroadcastEventType {
	return splitType
}
