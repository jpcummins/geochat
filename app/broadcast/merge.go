package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const mergeType types.BroadcastEventType = "merge"

type merge struct {
	Zone  interface{} `json:"zone"`
	Left  interface{} `json:"left"`
	Right interface{} `json:"right"`
}

func Merge(zone types.Zone, left types.Zone, right types.Zone) *merge {
	return &merge{
		Zone:  zone.BroadcastJSON(),
		Left:  left.BroadcastJSON(),
		Right: right.BroadcastJSON(),
	}
}

func (e *merge) Type() types.BroadcastEventType {
	return mergeType
}
