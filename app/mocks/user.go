package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type User struct {
	mock.Mock
}

func (m *User) ID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *User) Name() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *User) Location() types.LatLng {
	ret := m.Called()

	r0 := ret.Get(0).(types.LatLng)

	return r0
}
func (m *User) Zone() types.Zone {
	ret := m.Called()

	r0 := ret.Get(0).(types.Zone)

	return r0
}
func (m *User) SetZone(_a0 types.Zone) {
	m.Called(_a0)
}
func (m *User) Broadcast(_a0 types.BroadcastEventData) {
	m.Called(_a0)
}
func (m *User) ExecuteCommand(command string, args string) error {
	ret := m.Called(command, args)

	r0 := ret.Error(0)

	return r0
}
func (m *User) Connect() types.Connection {
	ret := m.Called()

	r0 := ret.Get(0).(types.Connection)

	return r0
}
func (m *User) Disconnect(_a0 types.Connection) {
	m.Called(_a0)
}

func (m *User) PubSubJSON() types.PubSubJSON {
	ret := m.Called()
	return ret.Get(0).(types.PubSubJSON)
}

func (m *User) Update(json types.PubSubJSON) error {
	ret := m.Called(json)
	return ret.Error(0)
}

func (m *User) BroadcastJSON() interface{} {
	ret := m.Called()
	return ret.Get(0)
}
