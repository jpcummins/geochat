package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type BroadcastEventData struct {
	mock.Mock
}

func (m *BroadcastEventData) Type() types.BroadcastEventType {
	ret := m.Called()

	r0 := ret.Get(0).(types.BroadcastEventType)

	return r0
}
func (m *BroadcastEventData) BeforeBroadcastToUser(_a0 types.User, _a1 types.BroadcastEvent) (bool, error) {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}
