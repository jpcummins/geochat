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
	SetLastSplit(time.Time)
	SetNextSplit(time.Time)
	SetLastMerge(time.Time)
	SetNextMerge(time.Time)

	// Broadcast
	Broadcast(BroadcastEventData) error

	// Pubsubs
	Join(User) error
	Leave(User) error
	Message(User, string) error
	Split() (map[string]Zone, error)
	Merge() error
}

type ZonePubSubJSON struct {
	ID        string    `json:"id"`
	UserIDs   []string  `json:"user_ids"`
	IsOpen    bool      `json:"is_open"`
	LastSplit time.Time `json:"last_split"`
	LastMerge time.Time `json:"last_merge"`
	NextSplit time.Time `json:"next_split"`
	NextMerge time.Time `json:"next_merge"`
}

func (psZone *ZonePubSubJSON) Type() PubSubDataType {
	return "zone"
}

type ZoneBroadcastJSON struct {
	ID        string                        `json:"id"`
	Users     map[string]*UserBroadcastJSON `json:"users"`
	SouthWest *LatLngJSON                   `json:"sw"`
	NorthEast *LatLngJSON                   `json:"ne"`
	NextSplit time.Time                     `json:"next_split"`
	NextMerge time.Time                     `json:"next_merge"`
}
