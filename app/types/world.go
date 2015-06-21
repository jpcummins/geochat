package types

type World interface {
	ID() string
	Users() Users
	Zones() Zones

	GetOrCreateZone(string) (Zone, error)
	GetOrCreateZoneForUser(User) (Zone, error)
	Publish(EventData) error
}
