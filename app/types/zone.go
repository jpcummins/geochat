package types

import "encoding/json"

type Zone interface {
	json.Marshaler

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
