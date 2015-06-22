package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type DB struct {
	mock.Mock
}

func (m *DB) User(userID string, worldID string) (*types.ServerUserJSON, error) {
	ret := m.Called(userID, worldID)

	var r0 *types.ServerUserJSON
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.ServerUserJSON)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SaveUser(json types.ServerJSON) error {
	ret := m.Called(json)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) Zone(zoneID string, worldID string) (*types.ServerZoneJSON, error) {
	ret := m.Called(zoneID, worldID)

	var r0 *types.ServerZoneJSON
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.ServerZoneJSON)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SaveZone(json types.ServerJSON) error {
	ret := m.Called(json)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) World(worldID string) (*types.ServerWorldJSON, error) {
	ret := m.Called(worldID)

	var r0 *types.ServerWorldJSON
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.ServerWorldJSON)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SaveWorld(_a0 types.ServerJSON) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
