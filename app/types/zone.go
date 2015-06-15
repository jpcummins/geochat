package types

import "encoding/json"

type Zone interface {
	json.Marshaler

	ID() string
	SouthWest() LatLng
	NorthEast() LatLng
	Geohash() string
	From() byte
	To() byte
	Parent() Zone
	Left() Zone
	Right() Zone
	MaxUsers() int
	Count() int
	IsOpen() bool
	SetIsOpen(bool)
	Broadcast(Event)
	AddUser(User)
	RemoveUser(string)
}
