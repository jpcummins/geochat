package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Event struct {
	mock.Mock
}

func (m *Event) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *Event) Type() string {
	args := m.Called()
	return args.String(0)
}

func (m *Event) Zone() types.Zone {
	args := m.Called()
	return args.Get(0).(types.Zone)
}

func (m *Event) Data() types.EventData {
	args := m.Called()
	return args.Get(0).(types.EventData)
}
