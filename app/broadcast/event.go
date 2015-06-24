package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

type Event struct {
	EventID   string                   `json:"id"`
	EventType types.BroadcastEventType `json:"type"`
	EventData types.BroadcastEventData `json:"data"`
}

func NewEvent(id string, data types.BroadcastEventData) *Event {
	return &Event{
		EventID:   id,
		EventType: data.Type(),
		EventData: data,
	}
}

func (e *Event) ID() string {
	return e.EventID
}

func (e *Event) Type() types.BroadcastEventType {
	return e.EventType
}

func (e *Event) Data() types.BroadcastEventData {
	return e.EventData
}
