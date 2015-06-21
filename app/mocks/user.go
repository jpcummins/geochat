package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *User) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *User) Location() types.LatLng {
	args := m.Called()
	return args.Get(0).(types.LatLng)
}

func (m *User) Zone() types.Zone {
	ret := m.Called()

	r0 := ret.Get(0).(types.Zone)

	return r0
}

func (m *User) Broadcast(e types.ClientEvent) {
	m.Called(e)
}

func (m *User) Connect() types.Connection {
	ret := m.Called()

	r0 := ret.Get(0).(types.Connection)

	return r0
}
func (m *User) Disconnect(_a0 types.Connection) {
	m.Called(_a0)
}
