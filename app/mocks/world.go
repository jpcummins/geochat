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
func (m *World) Users() types.Users {
	ret := m.Called()

	r0 := ret.Get(0).(types.Users)

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
func (m *World) NewUser(id string, name string, lat float64, lng float64) (types.User, error) {
	ret := m.Called(id, name, lat, lng)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *World) Publish(_a0 types.EventData) (types.Event, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Event)
	r1 := ret.Error(1)

	return r0, r1
}
