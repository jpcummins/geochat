package types

type DB interface {
	GetUser(string) (User, error)
	SetUser(User) error

	GetZone(zoneID string, worldID string) (Zone, error)
	SetZone(zoneID Zone, worldID string) error

	GetWorld(string) (World, error)
	SetWorld(World) error
}
