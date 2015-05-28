package chat

import (
	"encoding/json"
	"errors"
)

type messageJSON struct {
	SubscriptionID string `json:"subscription"`
	Text           string `json:"text"`
}

type Message struct {
	Subscriber *Subscription
	Text       string
}

func (m *Message) Type() string {
	return "message"
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var js messageJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	subscription, found := Subscribers.Get(js.SubscriptionID)
	if !found {
		panic(errors.New("Unknown subscription: " + js.SubscriptionID))
	}

	m.Subscriber = subscription
	m.Text = js.Text
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	messageJSON := &messageJSON{
		SubscriptionID: m.Subscriber.GetID(),
		Text:           m.Text,
	}
	return json.Marshal(messageJSON)
}

func (m *Message) OnReceive(e *Event) error {
	zone := m.Subscriber.GetZone()
	zone.broadcastEvent(e)
	zone.archiveEvent(e)
	return nil
}
