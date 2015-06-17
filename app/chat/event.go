package chat

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/events"
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
	data  types.EventData
	world types.World
}

func newEvent(id string, worldID string, data types.EventData) (*Event, error) {
	world, err := chat.World(worldID)

	if err != nil {
		return nil, err
	}

	if world == nil {
		return nil, errors.New("Unable to find world: " + worldID)
	}

	return &Event{
		eventJSON: &eventJSON{
			ID:      id,
			Type:    data.Type(),
			WorldID: worldID,
		},
		data:  data,
		world: world,
	}, nil
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

func (e *Event) World() types.World {
	return e.world
}

func (e *Event) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.eventJSON); err != nil {
		return err
	}

	world, err := chat.World(e.eventJSON.WorldID)
	if err != nil {
		return err
	}

	if world == nil {
		return errors.New("Unable to find world: " + e.eventJSON.WorldID)
	}

	e.world = world

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
