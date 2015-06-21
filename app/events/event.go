package events

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type eventJSON struct {
	ID      string                `json:"id"`
	Type    types.ServerEventType `json:"type"`
	WorldID string                `json:"world_id"`
	Data    json.RawMessage       `json:"data,omitempty"`
}

type ServerEvent struct {
	*eventJSON
	world types.World
	data  types.ServerEventData
}

func newEvent(id string, world types.World, data types.ServerEventData) *ServerEvent {
	return &ServerEvent{
		eventJSON: &eventJSON{
			ID:      id,
			Type:    data.Type(),
			WorldID: world.ID(),
		},
		world: world,
		data:  data,
	}
}

func (e *ServerEvent) ID() string {
	return e.eventJSON.ID
}

func (e *ServerEvent) Type() types.ServerEventType {
	return e.eventJSON.Type
}

func (e *ServerEvent) WorldID() string {
	return e.eventJSON.WorldID
}

func (e *ServerEvent) World() types.World {
	return e.world
}

func (e *ServerEvent) SetWorld(world types.World) {
	e.world = world
}

func (e *ServerEvent) Data() types.ServerEventData {
	return e.data
}

func (e *ServerEvent) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.eventJSON); err != nil {
		return err
	}

	switch e.Type() {
	case MessageServerEvent:
		e.data = &Message{}
	case JoinSeverEvent:
		e.data = &Join{}
	default:
		return errors.New("Unable to unmarshal command: " + string(e.Type()))
	}

	return json.Unmarshal(e.eventJSON.Data, e.data)
}
