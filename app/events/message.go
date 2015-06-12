package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Message struct {
	types.BaseEventData
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

func (m *Message) Type() string {
	return "message"
}

func (m *Message) OnReceive(e types.Event) error {
	e.Zone().Broadcast(e)
	return nil
}
