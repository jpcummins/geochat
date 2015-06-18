package types

type Cache interface {
	User(id string) (User, error)
	SetUser(user User) error
	UpdateUser(id string) (User, error)

	Zone(id string) (Zone, error)
	SetZone(zone Zone) error
	UpdateZone(id string) (Zone, error)

	World(id string) (World, error)
	SetWorld(world World) error
}
