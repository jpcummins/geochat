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

func (f *Events) NewServerEvent(id string, data types.ServerEventData) types.ServerEvent {
	return newServerEvent(id, f.world, data)
}

func (f *Events) NewClientEvent(id string, data types.ClientEventData) types.ClientEvent {
	return newClientEvent(id, data)
}
