package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Chat struct {
	mock.Mock
}

func (m *Chat) PubSub() types.PubSub {
	ret := m.Called()

	r0 := ret.Get(0).(types.PubSub)

	return r0
}
func (m *Chat) Cache() types.Cache {
	ret := m.Called()

	r0 := ret.Get(0).(types.Cache)

	return r0
}
func (m *Chat) Events() types.EventFactory {
	ret := m.Called()

	r0 := ret.Get(0).(types.EventFactory)

	return r0
}
func (m *Chat) World(id string) (types.World, error) {
	ret := m.Called(id)

	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}

	r0 := ret.Get(0).(types.World)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Chat) SetWorld(_a0 types.World) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *Chat) GetOrCreateWorld(id string) (types.World, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(types.World)
	r1 := ret.Error(1)

	return r0, r1
}
