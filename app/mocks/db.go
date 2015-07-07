package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type DB struct {
	mock.Mock
}

func (m *DB) User(userID string, worldID string) (*types.UserPubSubJSON, error) {
	ret := m.Called(userID, worldID)

	var r0 *types.UserPubSubJSON
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.UserPubSubJSON)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SaveUser(json *types.UserPubSubJSON, worldID string) error {
	ret := m.Called(json, worldID)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) Zone(zoneID string, worldID string) (*types.ZonePubSubJSON, error) {
	ret := m.Called(zoneID, worldID)

	var r0 *types.ZonePubSubJSON
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.ZonePubSubJSON)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SaveZone(json *types.ZonePubSubJSON, worldID string) error {
	ret := m.Called(json, worldID)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) BulkSaveUsersAndZones(users []*types.User, zones []*types.Zone, worldID string) error {
	ret := m.Called(users, zones, worldID)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) World(worldID string) (*types.WorldPubSubJSON, error) {
	ret := m.Called(worldID)

	var r0 *types.WorldPubSubJSON
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.WorldPubSubJSON)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SaveWorld(json *types.WorldPubSubJSON) error {
	ret := m.Called(json)

	r0 := ret.Error(0)

	return r0
}
