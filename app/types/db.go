package types

type DB interface {
	GetUser(string, User) (bool, error)
	SetUser(User) error

	GetZone(zoneID string, world World, zone Zone) (bool, error)
	SetZone(zoneID Zone, world World) error

	GetWorld(string, World) (bool, error)
	SetWorld(World) error
}
