package pubsub

import (
	"github.com/jpcummins/geochat/app/types"
)

const splitType types.PubSubEventType = "split"

type split struct {
	UserID string `json:"user_id"`
	ZoneID string `json:"zone_id"`
}

func Split(zone types.Zone, user types.User) (*split, error) {
	m := &split{
		UserID: user.ID(),
		ZoneID: zone.ID(),
	}
	return m, nil
}

func (m *split) Type() types.PubSubEventType {
	return splitType
}

func (m *split) BeforePublish(e types.PubSubEvent) error {
	return nil
}

func (m *split) OnReceive(e types.PubSubEvent) error {
	zone := e.World().Zones().FromCache(m.ZoneID)
	if zone == nil {
		return nil
	}
	return nil
}
