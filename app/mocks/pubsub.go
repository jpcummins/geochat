package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type PubSub struct {
	mock.Mock
}

func (m *PubSub) Publish(_a0 types.Event) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *PubSub) Subscribe() <-chan types.Event {
	ret := m.Called()

	var r0 chan types.Event
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(chan types.Event)
	}

	return r0
}
