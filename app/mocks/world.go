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
func (m *World) MaxUsers() int {
	ret := m.Called()

	r0 := ret.Get(0).(int)

	return r0
}
func (m *World) MinUsers() int {
	ret := m.Called()

	r0 := ret.Get(0).(int)

	return r0
}
func (m *World) Zones() types.Zones {
	ret := m.Called()

	r0 := ret.Get(0).(types.Zones)

	return r0
}

func (m *World) DB() types.DB {
	ret := m.Called()

	return ret.Get(0).(types.DB)
}

func (m *World) GetOrCreateZone(_a0 string) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *World) FindOpenZone(_a0 types.Zone, _a1 types.User) (types.Zone, error) {
	ret := m.Called(_a0, _a1)

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
func (m *World) Publish(_a0 types.PubSubEventData) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}

func (m *World) PubSubJSON() types.PubSubJSON {
	ret := m.Called()
	return ret.Get(0).(types.PubSubJSON)
}

func (m *World) Update(json types.PubSubJSON) error {
	ret := m.Called(json)
	return ret.Error(0)
}

func (m *World) BroadcastJSON() interface{} {
	ret := m.Called()
	return ret.Get(0)
}
func (m *World) Zone() types.Zone {
	ret := m.Called()

	r0 := ret.Get(0).(types.Zone)

	return r0
}
