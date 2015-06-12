package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) User(id string) (types.User, error) {
	args := m.Called(id)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockCache) SetUser(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockCache) Zone(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}

func (m *MockCache) SetZone(zone types.Zone) error {
	args := m.Called(zone)
	return args.Error(0)
}

func (m *MockCache) GetZoneForUser(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}
