package types

import "time"

type Zone interface {
	PubSubSerializable
	BroadcastSerializable

	ID() string
	World() World
	SouthWest() LatLng
	NorthEast() LatLng
	Geohash() string
	From() string
	To() string
	ParentZoneID() string
	LeftZoneID() string
	RightZoneID() string
	Count() int
	UserIDs() []string
	IsOpen() bool
	SetIsOpen(bool)
	AddUser(User)
	RemoveUser(string)

	// Broadcast
	Broadcast(BroadcastEventData) error

	Join(User) error
	Leave(User) error
	Message(User, string) error
	Split() (map[string]Zone, error)
	Merge() error
}

type ZonePubSubJSON struct {
	ID           string    `json:"id"`
	UserIDs      []string  `json:"user_ids"`
	IsOpen       bool      `json:"is_open"`
	LastModified time.Time `json:"last_modified"`
}

func (psZone *ZonePubSubJSON) Type() PubSubDataType {
	return "zone"
}

type ZoneBroadcastJSON struct {
	ID        string                        `json:"id"`
	Users     map[string]*UserBroadcastJSON `json:"users"`
	SouthWest *LatLngJSON                   `json:"sw"`
	NorthEast *LatLngJSON                   `json:"ne"`
}
