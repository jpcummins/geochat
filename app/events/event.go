package events

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type serverEventJSON struct {
	ID   string                `json:"id"`
	Type types.ServerEventType `json:"type"`
	Data types.ServerEventData `json:"data,omitempty"`
}

type ServerEvent struct {
	*serverEventJSON
	world types.World
}

func newServerEvent(id string, world types.World, data types.ServerEventData) *ServerEvent {
	return &ServerEvent{
		serverEventJSON: &serverEventJSON{
			ID:   id,
			Type: data.Type(),
			Data: data,
		},
		world: world,
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
	return e.serverEventJSON.Data
}

func (e *ServerEvent) UnmarshalJSON(b []byte) error {
	type AnonEvent struct {
		ID   string                `json:"id"`
		Type types.ServerEventType `json:"type"`
		Data json.RawMessage       `json:"data"`
	}
	var ae AnonEvent

	println("AA")

	if err := json.Unmarshal(b, &ae); err != nil {
		return err
	}

	println("BB")

	switch e.Type() {
	case MessageServerEvent:
		e.serverEventJSON.Data = &Message{}
	case JoinSeverEvent:
		e.serverEventJSON.Data = &serverJoin{}
	default:
		return errors.New("Unable to unmarshal command: " + string(e.Type()))
	}

	println("CC")

	e.serverEventJSON.ID = ae.ID
	e.serverEventJSON.Type = ae.Type
	return json.Unmarshal(ae.Data, e.Data)
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
