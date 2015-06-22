package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type World struct {
	mock.Mock
}

func (m *World) ID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *World) Zones() types.Zones {
	ret := m.Called()

	r0 := ret.Get(0).(types.Zones)

	return r0
}
func (m *World) GetOrCreateZone(_a0 string) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *World) FindOpenZone(_a0 types.User) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *World) Users() types.Users {
	ret := m.Called()

	r0 := ret.Get(0).(types.Users)

	return r0
}
func (m *World) NewUser(id string, name string, lat float64, lng float64) (types.User, error) {
	ret := m.Called(id, name, lat, lng)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *World) NewServerEvent(_a0 types.ServerEventData) types.ServerEvent {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.ServerEvent)

	return r0
}
func (m *World) NewClientEvent(_a0 types.ClientEventData) types.ClientEvent {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.ClientEvent)

	return r0
}
func (m *World) Publish(_a0 types.ServerEvent) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}

func (m *World) ClientJSON() types.ClientJSON {
	ret := m.Called()

	r0 := ret.Get(0).(types.ClientJSON)

	return r0
}
func (m *World) ServerJSON() types.ServerJSON {
	ret := m.Called()

	r0 := ret.Get(0).(types.ServerJSON)

	return r0
}
func (m *World) Update(_a0 types.ServerJSON) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
