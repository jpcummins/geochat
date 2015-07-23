package pubsub

import (
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const mergeType types.PubSubEventType = "merge"

type merge struct {
	Parent *types.ZonePubSubJSON `json:"parent"`
	Left   *types.ZonePubSubJSON `json:"left"`
	Right  *types.ZonePubSubJSON `json:"right"`
}

func Merge(parent types.Zone, left types.Zone, right types.Zone) (*merge, error) {
	s := &merge{
		Parent: parent.PubSubJSON().(*types.ZonePubSubJSON),
		Left:   left.PubSubJSON().(*types.ZonePubSubJSON),
		Right:  right.PubSubJSON().(*types.ZonePubSubJSON),
	}
	return s, nil
}

func (m *merge) Type() types.PubSubEventType {
	return mergeType
}

func (m *merge) BeforePublish(e types.PubSubEvent) error {
	return nil
}

func (m *merge) OnReceive(e types.PubSubEvent) error {
	parent := e.World().Zones().FromCache(m.Parent.ID)
	if parent == nil {
		return nil
	}

	parent.Update(m.Parent)

	left := e.World().Zones().FromCache(m.Left.ID)
	if left != nil {
		left.Update(m.Left)
	}

	right := e.World().Zones().FromCache(m.Right.ID)
	if right != nil {
		right.Update(m.Right)
	}

	parent.Broadcast(broadcast.Merge(parent, left, right))
	return nil
}
