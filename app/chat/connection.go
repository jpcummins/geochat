package chat

import (
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

type Connection struct {
	events chan types.BroadcastEvent
	user   types.User
}

func newConnection(user types.User) *Connection {
	return &Connection{
		events: make(chan types.BroadcastEvent, 10),
		user:   user,
	}
}

func (c *Connection) Events() <-chan types.BroadcastEvent {
	return c.events
}

func (c *Connection) Ping() {
	event := broadcast.NewEvent(generateEventID(), broadcast.Ping())
	c.events <- event
}
