package events

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type messageJSON struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type Message struct {
	*messageJSON
	user types.User
}

func NewMessage(user types.User, message string) (*Message, error) {
	m := &Message{
		messageJSON: &messageJSON{
			UserID:  user.ID(),
			Message: message,
		},
		user: user,
	}
	return m, nil
}

func (m *Message) Type() string {
	return "message"
}

func (m *Message) BeforePublish(e types.Event) error {

	if m.user.Zone() == nil {
		return errors.New("User is not associated with a zone")
	}

	return nil
}

func (m *Message) OnReceive(e types.Event) error {
	user, err := e.World().Users().User(m.UserID)
	if err != nil {
		return errors.New("Unable to lookup user " + m.UserID)
	}

	zone := user.Zone()
	if zone == nil {
		return errors.New("User is not associated with a zone")
	}

	zone.Broadcast(e)
	return nil
}
