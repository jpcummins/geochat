package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const mergeType types.BroadcastEventType = "merge"

type merge struct {
	Zone    interface{} `json:"zone"`
	LeftID  interface{} `json:"leftID"`
	RightID interface{} `json:"rightID"`
}

func Merge(zone types.Zone, leftID string, rightID string) *merge {
	return &merge{
		Zone:    zone.BroadcastJSON(),
		LeftID:  leftID,
		RightID: rightID,
	}
}

func (e *merge) Type() types.BroadcastEventType {
	return mergeType
}
