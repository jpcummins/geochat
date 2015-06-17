package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Event struct {
	mock.Mock
}

func (m *Event) ID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Event) World() types.World {
	ret := m.Called()

	r0 := ret.Get(0).(types.World)

	return r0
}
func (m *Event) Type() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Event) Data() types.EventData {
	ret := m.Called()

	r0 := ret.Get(0).(types.EventData)

	return r0
}

type EventData struct {
	mock.Mock
}

func (m *EventData) Type() string {
	ret := m.Called()
	r0 := ret.Get(0).(string)
	return r0
}

func (m *EventData) BeforePublish(_a0 types.Event) error {
	ret := m.Called(_a0)
	r0 := ret.Error(0)
	return r0
}

func (m *EventData) OnReceive(_a0 types.Event) error {
	ret := m.Called(_a0)
	r0 := ret.Error(0)
	return r0
}
