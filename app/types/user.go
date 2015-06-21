package types

type User interface {
	ID() string
	Name() string
	Location() LatLng
	Zone() Zone
	Broadcast(ClientEvent)
	Connect() Connection
	Disconnect(Connection)
}
