package types

type DB interface {
	GetUser(string) (User, error)
	SetUser(User) error

	GetZone(string) (Zone, error)
	SetZone(Zone) error

	GetWorld(string) (World, error)
	SetWorld(World) error
}
