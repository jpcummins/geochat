package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Chat struct {
	mock.Mock
}

func (m *Chat) DB() types.DB {
	ret := m.Called()

	r0 := ret.Get(0).(types.DB)

	return r0
}
func (m *Chat) PubSub() types.PubSub {
	ret := m.Called()

	r0 := ret.Get(0).(types.PubSub)

	return r0
}
func (m *Chat) Events() types.EventFactory {
	ret := m.Called()

	r0 := ret.Get(0).(types.EventFactory)

	return r0
}
func (m *Chat) World() types.World {
	ret := m.Called()

	r0 := ret.Get(0).(types.World)

	return r0
}
