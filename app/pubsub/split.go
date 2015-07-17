package pubsub

import (
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const splitType types.PubSubEventType = "split"

type split struct {
	Parent *types.ZonePubSubJSON `json:"parent"`
	Zones  []string              `json:"zones"`
}

func Split(parent types.Zone, zones map[string]types.Zone) (*split, error) {
	s := &split{}
	s.Parent = parent.PubSubJSON().(*types.ZonePubSubJSON)
	s.Zones = make([]string, 0, len(zones))

	for id := range zones {
		println(id)
		s.Zones = append(s.Zones, id)
	}
	return s, nil
}

func (m *split) Type() types.PubSubEventType {
	return splitType
}

func (m *split) BeforePublish(e types.PubSubEvent) error {
	return nil
}

func (m *split) OnReceive(e types.PubSubEvent) error {
	parent := e.World().Zones().FromCache(m.Parent.ID)
	if parent == nil {
		return nil
	}

	parent.Update(m.Parent)

	for _, id := range m.Zones {
		zone, err := e.World().Zones().FromDB(id)
		if err != nil {
			return err
		}
		zone.Broadcast(broadcast.Split(parent, zone))
	}

	return nil
}
