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
func (m *World) GetOrCreateZoneForUser(_a0 types.User) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *World) Publish(_a0 types.EventData) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
