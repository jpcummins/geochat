package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Events struct {
	mock.Mock
}

func (m *Events) NewServerEvent(_a0 string, _a1 types.ServerEventData) types.ServerEvent {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(types.ServerEvent)

	return r0
}
func (m *Events) NewClientEvent(_a0 string, _a1 types.ClientEventData) types.ClientEvent {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(types.ClientEvent)

	return r0
}
