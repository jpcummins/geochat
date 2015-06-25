package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type ServerEvent struct {
	mock.Mock
}

func (m *ServerEvent) ID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *ServerEvent) Type() types.ServerEventType {
	ret := m.Called()

	r0 := ret.Get(0).(types.ServerEventType)

	return r0
}
func (m *ServerEvent) WorldID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *ServerEvent) World() types.World {
	ret := m.Called()

	r0 := ret.Get(0).(types.World)

	return r0
}
func (m *ServerEvent) SetWorld(_a0 types.World) {
	m.Called(_a0)
}
func (m *ServerEvent) Data() types.ServerEventData {
	ret := m.Called()

	r0 := ret.Get(0).(types.ServerEventData)

	return r0
}

type ServerEventData struct {
	mock.Mock
}

func (m *ServerEventData) Type() types.ServerEventType {
	ret := m.Called()

	r0 := ret.Get(0).(types.ServerEventType)

	return r0
}
func (m *ServerEventData) BeforePublish(_a0 types.ServerEvent) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *ServerEventData) OnReceive(_a0 types.ServerEvent) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
