package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Cache struct {
	mock.Mock
}

func (m *Cache) User(id string, user types.User) (bool, error) {
	ret := m.Called(id, user)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) SetUser(user types.User) error {
	ret := m.Called(user)

	r0 := ret.Error(0)

	return r0
}
func (m *Cache) UpdateUser(id string, user types.User) (bool, error) {
	ret := m.Called(id, user)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) Zone(zoneID string, worldID string, zone types.Zone) (bool, error) {
	ret := m.Called(zoneID, worldID, zone)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) SetZone(zone types.Zone, worldID string) error {
	ret := m.Called(zone, worldID)

	r0 := ret.Error(0)

	return r0
}
func (m *Cache) UpdateZone(id string, worldID string, zone types.Zone) (bool, error) {
	ret := m.Called(id, worldID, zone)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) World(id string, world types.World) (bool, error) {
	ret := m.Called(id, world)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) SetWorld(world types.World) error {
	ret := m.Called(world)

	r0 := ret.Error(0)

	return r0
}
