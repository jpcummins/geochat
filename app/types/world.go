package types

type World interface {
	ID() string

	Zone(string) (Zone, error)
	SetZone(Zone) error
	UpdateZone(string) (Zone, error)
	GetOrCreateZone(string) (Zone, error)
	GetOrCreateZoneForUser(User) (Zone, error)

	User(string) (User, error)
	SetUser(User) error
	UpdateUser(string) (User, error)

	Publish(EventData) error
}
