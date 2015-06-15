package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type World struct {
	mock.Mock
}

func (m *World) Factory() types.Factory {
	args := m.Called()
	return args.Get(0).(types.Factory)
}

func (m *World) GetOrCreateZone(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}

func (m *World) MaxUsersForNewZones() int {
	args := m.Called()
	return args.Int(0)
}
