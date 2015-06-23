package events

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type serverEventJSON struct {
	ID   string                `json:"id"`
	Type types.ServerEventType `json:"type"`
	Data json.RawMessage       `json:"data,omitempty"`
}

type ServerEvent struct {
	*serverEventJSON
	world types.World
	data  types.ServerEventData
}

func newServerEvent(id string, world types.World, data types.ServerEventData) *ServerEvent {
	return &ServerEvent{
		serverEventJSON: &serverEventJSON{
			ID:   id,
			Type: data.Type(),
		},
		world: world,
		data:  data,
	}
}

func (e *ServerEvent) ID() string {
	return e.serverEventJSON.ID
}

func (e *ServerEvent) Type() types.ServerEventType {
	return e.serverEventJSON.Type
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
	if err := json.Unmarshal(b, &e.serverEventJSON); err != nil {
		return err
	}

	switch e.Type() {
	case MessageServerEvent:
		e.data = &Message{}
	case JoinSeverEvent:
		e.data = &serverJoin{}
	default:
		return errors.New("Unable to unmarshal command: " + string(e.Type()))
	}

	return json.Unmarshal(e.serverEventJSON.Data, e.data)
}

type clientEventJSON struct {
	ID   string                `json:"id"`
	Type types.ClientEventType `json:"type"`
	Data types.ClientEventData `json:"data,omitempty"`
}

type ClientEvent struct {
	*clientEventJSON
}

func newClientEvent(id string, data types.ClientEventData) *ClientEvent {
	return &ClientEvent{
		clientEventJSON: &clientEventJSON{
			ID:   id,
			Type: data.Type(),
			Data: data,
		},
	}
}

func (c *ClientEvent) ID() string {
	return c.clientEventJSON.ID
}

func (c *ClientEvent) Type() types.ClientEventType {
	return c.clientEventJSON.Type
}

func (c *ClientEvent) Data() types.ClientEventData {
	return c.clientEventJSON.Data
}
