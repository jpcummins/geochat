package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type World struct {
	mock.Mock
}

func (m *World) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *World) GetOrCreateZone(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}

func (m *World) GetOrCreateZoneForUser(user types.User) (types.Zone, error) {
	args := m.Called(user)
	return args.Get(0).(types.Zone), args.Error(1)
}

func (m *World) Publish(event types.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
