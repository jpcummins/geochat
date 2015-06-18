package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Cache struct {
	mock.Mock
}

func (m *Cache) User(id string) (types.User, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) SetUser(user types.User) error {
	ret := m.Called(user)

	r0 := ret.Error(0)

	return r0
}
func (m *Cache) UpdateUser(id string) (types.User, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) Zone(id string) (types.Zone, error) {
	ret := m.Called(id)

	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) SetZone(zone types.Zone) error {
	ret := m.Called(zone)

	r0 := ret.Error(0)

	return r0
}
func (m *Cache) UpdateZone(id string) (types.Zone, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) World(id string) (types.World, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(types.World)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Cache) SetWorld(world types.World) error {
	ret := m.Called(world)

	r0 := ret.Error(0)

	return r0
}
