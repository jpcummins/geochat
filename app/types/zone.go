package types

type Zone interface {
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

	ClientJSON() *ClientZoneJSON
	ServerJSON() *ServerZoneJSON
	Update(*ServerZoneJSON) error

	// Events
	Join(User) (ClientEvent, error)
	Message(User, string) (ClientEvent, error)
}

type ServerZoneJSON struct {
	ID       string   `json:"id"`
	UserIDs  []string `json:"user_ids"`
	IsOpen   bool     `json:"is_open"`
	MaxUsers int      `json:"max_users"`
}

type ClientZoneJSON struct {
	ID    string            `json:"id"`
	Users []*ClientUserJSON `json:"users"`
}
