package types

type Cache interface {
	User(id string, user User) (bool, error)
	SetUser(user User) error
	UpdateUser(id string, user User) (bool, error)

	Zone(zoneID string, worldID string, zone Zone) (bool, error)
	SetZone(zone Zone, worldID string) error
	UpdateZone(id string, worldID string, zone Zone) (bool, error)

	World(id string, world World) (bool, error)
	SetWorld(world World) error
}
