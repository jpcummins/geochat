package types

import "encoding/json"

type Zone interface {
	json.Marshaler

	ID() string
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
	Broadcast(Event)
	AddUser(User)
	RemoveUser(string)

	// Events
	Join(User) (Event, error)
	Message(User, string) (Event, error)
}
