package chat

import (
	"encoding/json"
	"errors"
	"time"
)

type EventData interface {
	Type() string
}

type Event struct {
	Type      string    `json:"type"`
	Timestamp int       `json:"timestamp"`
	Data      EventData `json:"data,omitempty"`
}

func (e *Event) UnmarshalJSON(b []byte) error {

	type AnonEvent struct {
		Type      string          `json:"type"`
		Timestamp int             `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	}
	var ae AnonEvent

	if err := json.Unmarshal(b, &ae); err != nil {
		return err
	}

	switch ae.Type {
	case "message":
		e.Data = &Message{}
	case "join":
		e.Data = &Join{}
	case "leave":
		e.Data = &Leave{}
	case "online":
		e.Data = &Online{}
	case "offline":
		e.Data = &Offline{}
	default:
		return errors.New("Unable to unmarshal command: " + ae.Type)
	}

	e.Type = ae.Type
	e.Timestamp = ae.Timestamp

	return json.Unmarshal(ae.Data, e.Data)
}

func NewEvent(data EventData) *Event {
	return &Event{Type: data.Type(), Data: data, Timestamp: int(time.Now().Unix())}
}
