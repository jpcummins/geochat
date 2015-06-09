package chat

import (
	"encoding/json"
	"errors"
)

type messageJSON struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

type Message struct {
	User *User
	Text string
}

func (m *Message) Type() string {
	return "message"
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var js messageJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	user, found := UserCache.Get(js.UserID)
	if !found {
		panic(errors.New("Unknown user: " + js.UserID))
	}

	m.User = user
	m.Text = js.Text
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	messageJSON := &messageJSON{
		UserID: m.User.GetID(),
		Text:   m.Text,
	}
	return json.Marshal(messageJSON)
}

func (m *Message) OnReceive(e *Event) error {
	println("sending")
	zone := m.User.GetZone()
	zone.broadcastEvent(e)
	zone.archiveEvent(e)
	return nil
}
