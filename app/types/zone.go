package types

type Zone interface {
	EventJSON

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
	Broadcast(ClientEvent)
	AddUser(User)
	RemoveUser(string)

	// Events
	Join(User) (ClientEvent, error)
	Message(User, string) (ClientEvent, error)
}

type ServerZoneJSON struct {
	*BaseServerJSON
	UserIDs  []string `json:"user_ids"`
	IsOpen   bool     `json:"is_open"`
	MaxUsers int      `json:"max_users"`
}

type ClientZoneJSON struct {
	BaseClientJSON
	Users []*ClientUserJSON `json:"users"`
}
