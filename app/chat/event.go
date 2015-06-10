package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type EventData interface {
	Type() string
	OnReceive(*Event) error
}

type Event struct {
	Type string    `json:"type"`
	ID   string    `json:"id"`
	Data EventData `json:"data,omitempty"`
}

func (e *Event) UnmarshalJSON(b []byte) error {

	type AnonEvent struct {
		Type string          `json:"type"`
		ID   string          `json:"id"`
		Data json.RawMessage `json:"data"`
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
	case "split":
		e.Data = &Split{}
	default:
		return errors.New("Unable to unmarshal command: " + ae.Type)
	}

	e.Type = ae.Type
	e.ID = ae.ID

	return json.Unmarshal(ae.Data, e.Data)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewEvent(data EventData) *Event {
	return &Event{Type: data.Type(), Data: data, ID: fmt.Sprintf("%d%s", time.Now().Unix(), randomSequence(4))}
}
