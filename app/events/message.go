package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type messageJSON struct {
	UserID  string `json:"user_id"`
	ZoneID  string `json:"zone_id"`
	Message string `json:"message"`
}

type Message struct {
	*messageJSON
	user types.User
	zone types.Zone
}

func NewMessage(user types.User, zone types.Zone, message string) (*Message, error) {
	m := &Message{
		messageJSON: &messageJSON{
			UserID:  user.ID(),
			Message: message,
		},
		user: user,
		zone: zone,
	}
	return m, nil
}

func (m *Message) Type() string {
	return "message"
}

func (m *Message) BeforePublish(e types.Event) error {
	// Archive message
	return nil
}

func (m *Message) OnReceive(e types.Event) error {
	zone, err := e.World().Zones().Zone(m.messageJSON.ZoneID)
	if err != nil {
		return nil
	}

	zone.Broadcast(e)
	return nil
}
