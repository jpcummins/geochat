package types

type DB interface {
	Publish(Event) error
	Subscribe(Zone) <-chan Event

	GetUser(string) (User, error)
	SetUser(User) error

	GetZone(string) (Zone, error)
	SetZone(Zone) error
}
