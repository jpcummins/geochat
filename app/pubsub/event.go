package pubsub

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type Event struct {
	EventID   string                `json:"id"`
	WorldID   string                `json:"world_id"`
	EventType types.PubSubEventType `json:"type"`
	EventData types.PubSubEventData `json:"data,omitempty"`
	world     types.World
}

func NewEvent(id string, world types.World, data types.PubSubEventData) *Event {
	return &Event{
		EventID:   id,
		WorldID:   world.ID(),
		EventType: data.Type(),
		EventData: data,
		world:     world,
	}
}

func (e *Event) ID() string {
	return e.EventID
}

func (e *Event) Type() types.PubSubEventType {
	return e.EventType
}

func (e *Event) Data() types.PubSubEventData {
	return e.EventData
}

func (e *Event) World() types.World {
	return e.world
}

func (e *Event) SetWorld(world types.World) {
	e.world = world
}

func (e *Event) UnmarshalJSON(b []byte) error {
	type AnonEvent struct {
		ID      string                `json:"id"`
		WorldID string                `json:"world_id"`
		Type    types.PubSubEventType `json:"type"`
		Data    json.RawMessage       `json:"data"`
	}

	var ae AnonEvent
	if err := json.Unmarshal(b, &ae); err != nil {
		return err
	}

	switch ae.Type {
	case messageType:
		e.EventData = &message{}
	case joinType:
		e.EventData = &join{}
	case leaveType:
		e.EventData = &leave{}
	case splitType:
		e.EventData = &split{}
	case mergeType:
		e.EventData = &merge{}
	case worldType:
		e.EventData = &world{}
	case announcementType:
		e.EventData = &announcement{}
	default:
		return errors.New("Unable to unmarshal command: " + string(ae.Type))
	}

	e.EventID = ae.ID
	e.EventType = ae.Type
	e.WorldID = ae.WorldID
	return json.Unmarshal(ae.Data, e.EventData)
}
