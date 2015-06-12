package types

type User interface {
	ID() string
	Name() string
	Broadcast(Event)
	AddConnection(Connection)
	RemoveConnection(Connection)
}
