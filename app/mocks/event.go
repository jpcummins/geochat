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

func (m *Event) WorldID() string {
	args := m.Called()
	return args.String(0)
}

func (m *Event) Data() types.EventData {
	args := m.Called()
	return args.Get(0).(types.EventData)
}

func (m *Event) UnmarshalJSON(b []byte) error {
	args := m.Called(b)
	return args.Error(0)
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
