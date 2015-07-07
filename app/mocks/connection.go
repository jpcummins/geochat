package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Connection struct {
	mock.Mock
}

func (m *Connection) Events() chan types.PubSubEvent {
	args := m.Called()
	return args.Get(0).(chan types.PubSubEvent)
}

func (m *Connection) Ping() {
	m.Called()
}
