package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
	"math/rand"
	"time"
)

type eventJSON struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	WorldID string          `json:"world_id"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Event struct {
	*eventJSON
	data types.EventData
}

func newEvent(world types.World, data types.EventData) *Event {
	return &Event{
		eventJSON: &eventJSON{
			ID:      fmt.Sprintf("%d:%s", time.Now().UnixNano(), randomSequence(4)),
			Type:    data.Type(),
			WorldID: world.ID(),
		},
	}
}

func (e *Event) ID() string {
	return e.eventJSON.ID
}

func (e *Event) Type() string {
	return e.eventJSON.Type
}

func (e *Event) Data() types.EventData {
	return e.data
}

func (e *Event) WorldID() string {
	return e.eventJSON.WorldID
}

func (e *Event) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.eventJSON); err != nil {
		return err
	}

	switch e.Type() {
	case "message":
		e.data = &events.Message{}
	case "join":
		e.data = &events.Join{}
	case "leave":
		e.data = &events.Leave{}
	case "online":
		e.data = &events.Online{}
	case "offline":
		e.data = &events.Offline{}
	case "split":
		e.data = &events.Split{}
	default:
		return errors.New("Unable to unmarshal command: " + e.Type())
	}

	return json.Unmarshal(e.eventJSON.Data, e.data)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
