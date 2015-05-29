package chat

type Connection struct {
	user   *User
	Events chan *Event
}

func newConnection(user *User) *Connection {
	return &Connection{
		user:   user,
		Events: make(chan *Event, 10),
	}
}
