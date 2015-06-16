package types

type User interface {
	ID() string
	Name() string
	Location() LatLng
	Broadcast(Event)
	AddConnection(Connection)
	RemoveConnection(Connection)
}
