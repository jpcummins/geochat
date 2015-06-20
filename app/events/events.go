package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Events struct {
	world types.World
}

func NewEvents(world types.World) types.Events {
	return &Events{
		world: world,
	}
}

func (f *Events) New(id string, data types.EventData) (types.Event, error) {
	return newEvent(id, f.world, data)
}
