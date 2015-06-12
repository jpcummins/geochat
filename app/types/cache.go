package types

type Cache interface {
	User(id string) (User, error)
	SetUser(user User) error

	Zone(id string) (Zone, error)
	SetZone(zone Zone) error

	GetZoneForUser(id string) (Zone, error)
}
