package types

import (
	"encoding/json"
)

type Event interface {
	ID() string
	Type() string
	World() World
	Zone() Zone
	Data() EventData
	json.Marshaler
	json.Unmarshaler
}

type EventData interface {
	Type() string
	World() World
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
