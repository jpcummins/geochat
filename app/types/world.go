package types

type World interface {
	EventJSON

	ID() string

	Zones() Zones
	GetOrCreateZone(string) (Zone, error)
	FindOpenZone(User) (Zone, error)

	Users() Users
	NewUser(id string, name string, lat float64, lng float64) (User, error)

	NewServerEvent(ServerEventData) ServerEvent
	NewClientEvent(ClientEventData) ClientEvent
	Publish(ServerEvent) error
}

type ServerWorldJSON BaseServerJSON

type ClientWorldJSON BaseClientJSON
