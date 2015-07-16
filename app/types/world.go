package types

import "time"

type World interface {
	PubSubSerializable

	ID() string
	MaxUsers() int
	MinUsers() int
	SplitDelay() time.Duration
	DB() DB
	Zone() Zone

	Zones() Zones
	GetOrCreateZone(string) (Zone, error)
	FindOpenZone(Zone, User) (Zone, error)

	Users() Users
	NewUser(id string, name string, lat float64, lng float64) (User, error)

	Publish(PubSubEventData) error
}

type WorldPubSubJSON struct {
	ID         string        `json:"id"`
	MaxUsers   int           `json:"max_users"`
	MinUsers   int           `json:"min_users"`
	SplitDelay time.Duration `json:"split_delay"`
}

func (psWorld *WorldPubSubJSON) Type() PubSubDataType {
	return "world"
}
