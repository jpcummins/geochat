package types

type DB interface {
	Publish(event Event) error
	Subscribe(zone Zone) <-chan Event

	GetUser(id string) (User, error)
	SetUser(user User) error

	GetZone(id string) (Zone, error)
	SetZone(zone Zone) error
}
