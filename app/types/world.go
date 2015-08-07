package types

import "time"

type World interface {
	PubSubSerializable

	ID() string

	MaxUsers() int
	SetMaxUsers(int)

	MinUsers() int
	SetMinUsers(int)

	SplitDelay() time.Duration
	DB() DB
	Zone() Zone

	Zones() Zones
	GetOrCreateZone(string) (Zone, error)
	FindOpenZone(Zone, User) (Zone, error)
	Join(User) (Zone, error)

	Users() Users
	NewUser(string) (User, error)

	Publish(PubSubEventData) error
}

type WorldPubSubJSON struct {
	ID         string        `json:"id"`
	MaxUsers   int           `json:"maxUsers"`
	MinUsers   int           `json:"minUsers"`
	SplitDelay time.Duration `json:"splitDelay"`
}

func (psWorld *WorldPubSubJSON) Type() PubSubDataType {
	return "world"
}
