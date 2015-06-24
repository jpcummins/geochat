package pubsub

import (
	"errors"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const messageType types.PubSubEventType = "message"

type message struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	user    types.User
}

func Message(user types.User, text string) (*message, error) {
	m := &message{
		UserID:  user.ID(),
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
	var user types.User
	if user = e.World().Users().FromCache(m.UserID); user == nil {
		return errors.New("Unknown user")
	}
	user.Zone().Broadcast(broadcast.Message(m.UserID, m.Message))
	return nil
}
