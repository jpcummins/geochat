package pubsub

import (
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const messageType types.PubSubEventType = "message"

type message struct {
	UserID  string `json:"user_id"`
	ZoneID  string `json:"zone_id"`
	Message string `json:"message"`
	user    types.User
}

func Message(user types.User, zone types.Zone, text string) (*message, error) {
	m := &message{
		UserID:  user.ID(),
		ZoneID:  zone.ID(),
		Message: text,
		user:    user,
	}
	return m, nil
}

func (m *message) Type() types.PubSubEventType {
	return messageType
}

func (m *message) BeforePublish(e types.PubSubEvent) error {
	return nil
}

func (m *message) OnReceive(e types.PubSubEvent) error {
	zone, err := e.World().Zones().Zone(m.ZoneID)
	if err != nil {
		return err
	}

	zone.Broadcast(broadcast.Message(m.UserID, m.Message))
	return nil
}
