package types

type Cache interface {
	User(id string) (User, error)
	SetUser(user User) error
	UpdateUser(id string) (User, error)

	Zone(zoneID string, worldID string) (Zone, error)
	SetZone(zone Zone, worldID string) error
	UpdateZone(id string, worldID string) (Zone, error)

	World(id string) (World, error)
	SetWorld(world World) error
}
