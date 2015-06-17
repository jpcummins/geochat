package types

type World interface {
	ID() string

	Zone(string) (Zone, error)
	SetZone(Zone) error
	GetOrCreateZone(string) (Zone, error)
	GetOrCreateZoneForUser(User) (Zone, error)

	SetUser(User) error
	Publish(Event) error
}
