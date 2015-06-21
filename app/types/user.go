package types

type User interface {
	ID() string
	Name() string
	Location() LatLng
	Zone() Zone
	Broadcast(Event)
	Connect() Connection
	Disconnect(Connection)
}
