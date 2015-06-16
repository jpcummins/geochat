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
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Data types.EventData `json:"data,omitempty"`
}

type Event struct {
	*eventJSON
	world types.World
	zone  types.Zone
}

func (e *Event) ID() string {
	return e.eventJSON.ID
}

func (e *Event) Type() string {
	return e.eventJSON.Type
}

func (e *Event) World() types.World {
	return e.world
}

func (e *Event) Zone() types.Zone {
	return e.zone
}

func (e *Event) Data() types.EventData {
	return e.eventJSON.Data
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
		e.Data = &events.Message{}
	case "join":
		e.Data = &events.Join{}
	case "leave":
		e.Data = &events.Leave{}
	case "online":
		e.Data = &events.Online{}
	case "offline":
		e.Data = &events.Offline{}
	case "split":
		e.Data = &events.Split{}
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

func newEvent(data types.EventData) *Event {
	return &Event{
		Type: data.Type(),
		Data: data,
		ID:   fmt.Sprintf("%s:%d%s", data.World().ID(), time.Now().Unix(), randomSequence(4)),
	}
}
