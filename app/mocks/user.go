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

func (m *User) SetZone(_a0 types.Zone) {
	m.Called(_a0)
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

func (m *User) ClientJSON() types.ClientJSON {
	ret := m.Called()

	r0 := ret.Get(0).(types.ClientJSON)

	return r0
}
func (m *User) ServerJSON() types.ServerJSON {
	ret := m.Called()

	r0 := ret.Get(0).(types.ServerJSON)

	return r0
}
func (m *User) Update(_a0 types.ServerJSON) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
