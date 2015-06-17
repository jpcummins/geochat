package types

import (
	"encoding/json"
)

type Event interface {
	ID() string
	World() World
	Type() string
	Data() EventData
	json.Unmarshaler
}

type EventData interface {
	Type() string
	BeforePublish(Event) error
	OnReceive(Event) error
}

type BaseEventData struct{}

func (d *BaseEventData) BeforePublish(event Event) error {
	return nil
}

func (d *BaseEventData) OnReceive(event Event) error {
	return nil
}
