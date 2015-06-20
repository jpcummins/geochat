package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type DB struct {
	mock.Mock
}

func (m *DB) GetUser(_a0 string, _a1 types.User) (bool, error) {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetUser(_a0 types.User) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) GetZone(zoneID string, world types.World, zone types.Zone) (bool, error) {
	ret := m.Called(zoneID, world, zone)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetZone(zoneID types.Zone, world types.World) error {
	ret := m.Called(zoneID, world)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) GetWorld(_a0 string, _a1 types.World) (bool, error) {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetWorld(_a0 types.World) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
