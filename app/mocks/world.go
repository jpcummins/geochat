package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type World struct {
	mock.Mock
}

func (m *World) GetOrCreateZone(id string) (types.Zone, error) {
	args := m.Called(id)
	return args.Get(0).(types.Zone), args.Error(1)
}
