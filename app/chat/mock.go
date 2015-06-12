package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type mockUser struct {
	mock.Mock
}

func (m *mockUser) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockUser) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockUser) Broadcast(e types.Event) {
	m.Called(e)
}

func (m *mockUser) AddConnection(c types.Connection) {
	m.Called(c)
}

func (m *mockUser) RemoveConnection(connection types.Connection) {
	m.Called(connection)
}

type mockConnection struct {
	mock.Mock
}

func (m *mockConnection) Events() chan types.Event {
	args := m.Called()
	return args.Get(0).(chan types.Event)
}

type mockEvent struct {
	mock.Mock
}

func (m *mockEvent) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockEvent) Type() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockEvent) Zone() types.Zone {
	args := m.Called()
	return args.Get(0).(types.Zone)
}

func (m *mockEvent) Data() types.EventData {
	args := m.Called()
	return args.Get(0).(types.EventData)
}
