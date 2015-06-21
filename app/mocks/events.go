package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Events struct {
	mock.Mock
}

func (m *Events) New(_a0 string, _a1 types.EventData) (types.Event, error) {
	ret := m.Called(_a0, _a1)

	r0 := ret.Get(0).(types.Event)
	r1 := ret.Error(1)

	return r0, r1
}
