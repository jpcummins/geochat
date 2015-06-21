package types

type User interface {
	ID() string
	Name() string
	Location() LatLng
	Zone() Zone
	Broadcast(Event)
	AddConnection(Connection)
	RemoveConnection(Connection)
}
