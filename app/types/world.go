package types

type World interface {
	ID() string
	GetOrCreateZone(string) (Zone, error)
	GetOrCreateZoneForUser(User) (Zone, error)
	Publish(Event) error
}
