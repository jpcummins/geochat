package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Connection struct {
	events chan types.Event
}

func newConnection() *Connection {
	return &Connection{
		events: make(chan types.Event, 10),
	}
}

func (c *Connection) Events() chan types.Event {
	return c.events
}
