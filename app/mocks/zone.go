package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Zone struct {
	mock.Mock
}

func (z *Zone) ID() string {
	args := z.Called()
	return args.String(0)
}

func (z *Zone) SouthWest() types.LatLng {
	args := z.Called()
	return args.Get(0).(types.LatLng)
}

func (z *Zone) NorthEast() types.LatLng {
	args := z.Called()
	return args.Get(0).(types.LatLng)
}

func (z *Zone) Geohash() string {
	args := z.Called()
	return args.String(0)
}

func (z *Zone) From() string {
	args := z.Called()
	return args.String(0)
}

func (z *Zone) To() string {
	args := z.Called()
	return args.Get(0).(string)
}

func (z *Zone) ParentZoneID() string {
	args := z.Called()
	return args.String(0)
}

func (z *Zone) LeftZoneID() string {
	args := z.Called()
	return args.String(0)
}

func (z *Zone) RightZoneID() string {
	args := z.Called()
	return args.String(0)
}

func (z *Zone) MaxUsers() int {
	args := z.Called()
	return args.Int(0)
}

func (z *Zone) Count() int {
	args := z.Called()
	return args.Int(0)
}

func (z *Zone) IsOpen() bool {
	args := z.Called()
	return args.Bool(0)
}

func (z *Zone) SetIsOpen(isOpen bool) {
	z.Called(isOpen)
}

func (z *Zone) Broadcast(event types.Event) {
	z.Called(event)
}

func (z *Zone) AddUser(user types.User) {
	z.Called(user)
}

func (z *Zone) RemoveUser(id string) {
	z.Called(id)
}

func (z *Zone) MarshalJSON() ([]byte, error) {
	args := z.Called()
	return args.Get(0).([]byte), args.Error(1)
}
