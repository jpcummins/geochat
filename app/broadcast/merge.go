package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const mergeType types.BroadcastEventType = "merge"

type merge struct {
	Zone interface{} `json:"zone"`
}

func Merge(zone types.Zone) *merge {
	return &merge{
		Zone: zone.BroadcastJSON(),
	}
}

func (e *merge) Type() types.BroadcastEventType {
	return mergeType
}
