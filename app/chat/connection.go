package chat

import (
	"github.com/jpcummins/geochat/app/events"
	"github.com/jpcummins/geochat/app/types"
)

type Connection struct {
	events chan types.Event
	user   types.User
}

func newConnection(user types.User) *Connection {
	return &Connection{
		events: make(chan types.Event, 10),
		user:   user,
	}
}

func (c *Connection) Events() chan types.Event {
	return c.events
}

func (c *Connection) Ping() {
	c.events <- c.user.Zone().World().NewEvent(&events.Ping{})
}
