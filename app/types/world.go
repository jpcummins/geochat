package types

type World interface {
	ID() string
	Users() Users
	Zones() Zones

	GetOrCreateZone(string) (Zone, error)
	FindOpenZone(User) (Zone, error)
	Publish(EventData) error
}
