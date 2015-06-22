package types

type World interface {
	EventJSON

	ID() string
	MaxUsers() int

	Zones() Zones
	GetOrCreateZone(string) (Zone, error)
	FindOpenZone(User) (Zone, error)

	Users() Users
	NewUser(id string, name string, lat float64, lng float64) (User, error)

	NewServerEvent(ServerEventData) ServerEvent
	NewClientEvent(ClientEventData) ClientEvent
	Publish(ServerEvent) error
}

type ServerWorldJSON struct {
	*BaseServerJSON
	MaxUsers int `json:"max_users"`
}

type ClientWorldJSON BaseClientJSON
