package types

type World interface {
	ID() string

	Zones() Zones
	GetOrCreateZone(string) (Zone, error)
	FindOpenZone(User) (Zone, error)

	Users() Users
	NewUser(id string, name string, lat float64, lng float64) (User, error)

	Publish(EventData) error
}
