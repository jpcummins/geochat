package events

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type eventJSON struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	WorldID string          `json:"world_id"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Event struct {
	*eventJSON
	world types.World
	data  types.EventData
}

func newEvent(id string, world types.World, data types.EventData) (*Event, error) {
	return &Event{
		eventJSON: &eventJSON{
			ID:      id,
			Type:    data.Type(),
			WorldID: world.ID(),
		},
		world: world,
		data:  data,
	}, nil
}

func (e *Event) ID() string {
	return e.eventJSON.ID
}

func (e *Event) Type() string {
	return e.eventJSON.Type
}

func (e *Event) WorldID() string {
	return e.eventJSON.WorldID
}

func (e *Event) World() types.World {
	return e.world
}

func (e *Event) SetWorld(world types.World) {
	e.world = world
}

func (e *Event) Data() types.EventData {
	return e.data
}

func (e *Event) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.eventJSON); err != nil {
		return err
	}

	switch e.Type() {
	case "message":
		e.data = &Message{}
	case "join":
		e.data = &Join{}
	case "leave":
		e.data = &Leave{}
	case "online":
		e.data = &Online{}
	case "offline":
		e.data = &Offline{}
	case "split":
		e.data = &Split{}
	default:
		return errors.New("Unable to unmarshal command: " + e.Type())
	}

	return json.Unmarshal(e.eventJSON.Data, e.data)
}
