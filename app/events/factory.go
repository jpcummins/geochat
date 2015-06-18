package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type EventFactory struct {
	world types.World
}

func NewEventFactory(world types.World) types.EventFactory {
	return &EventFactory{
		world: world,
	}
}

func (f *EventFactory) New(id string, data types.EventData) (types.Event, error) {
	return newEvent(id, f.world, data)
}
