package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type MockUserCache struct {
	mock.Mock
}

func (m *MockUserCache) User(id string) (types.User, error) {
	args := m.Called(id)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockUserCache) SetUser(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserCache) Zone(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}

func (m *MockUserCache) SetZone(zone types.Zone) error {
	args := m.Called(zone)
	return args.Error(0)
}

func (m *MockUserCache) GetZoneForUser(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}
