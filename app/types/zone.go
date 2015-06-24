package types

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
	MaxUsers() int
	Count() int
	IsOpen() bool
	SetIsOpen(bool)
	AddUser(User)
	RemoveUser(string)

	// Broadcast
	Broadcast(BroadcastEventData)

	// Pubsubs
	Join(User) (BroadcastEvent, error)
	Message(User, string) (BroadcastEvent, error)
}

type ZonePubSubJSON struct {
	ID       string   `json:"id"`
	UserIDs  []string `json:"user_ids"`
	IsOpen   bool     `json:"is_open"`
	MaxUsers int      `json:"max_users"`
}

func (psZone *ZonePubSubJSON) Type() PubSubDataType {
	return "zone"
}

type ZoneBroadcastJSON struct {
	ID    string               `json:"id"`
	Users []*UserBroadcastJSON `json:"user_ids"`
}
