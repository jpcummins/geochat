package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Cache struct {
	mock.Mock
}

func (m *Cache) User(id string) (types.User, error) {
	args := m.Called(id)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *Cache) SetUser(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *Cache) Zone(id string) (types.Zone, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(types.Zone), args.Error(1)
}

func (m *Cache) SetZone(z types.Zone) error {
	args := m.Called(z)
	return args.Error(0)
}

func (m *Cache) World(id string) (types.World, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(types.World), args.Error(1)
}

func (m *Cache) SetWorld(z types.World) error {
	args := m.Called(z)
	return args.Error(0)
}
