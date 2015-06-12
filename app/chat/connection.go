package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Connection struct {
	user   *User
	events chan types.Event
}

func newConnection(user *User) *Connection {
	return &Connection{
		user:   user,
		events: make(chan types.Event, 10),
	}
}

func (c *Connection) Events() chan types.Event {
	return c.events
}
