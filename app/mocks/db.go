package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type DB struct {
	mock.Mock
}

func (m *DB) GetUser(_a0 string) (types.User, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetUser(_a0 types.User) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) GetZone(zoneID string, worldID string) (types.Zone, error) {
	ret := m.Called(zoneID, worldID)

	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetZone(zoneID types.Zone, worldID string) error {
	ret := m.Called(zoneID, worldID)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) GetWorld(_a0 string) (types.World, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.World)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetWorld(_a0 types.World) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
