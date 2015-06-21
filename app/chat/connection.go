package chat

import (
	// "github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
)

type Connection struct {
	events chan types.ClientEvent
	user   types.User
}

func newConnection(user types.User) *Connection {
	return &Connection{
		events: make(chan types.ClientEvent, 10),
		user:   user,
	}
}

func (c *Connection) Events() chan types.ClientEvent {
	return c.events
}

func (c *Connection) Ping() {
	// c.events <- c.user.Zone().World().NewEvent(&events.Ping{})
}
