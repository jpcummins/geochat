package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *User) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *User) Location() types.LatLng {
	args := m.Called()
	return args.Get(0).(types.LatLng)
}

func (m *User) Broadcast(e types.Event) {
	m.Called(e)
}

func (m *User) AddConnection(c types.Connection) {
	m.Called(c)
}

func (m *User) RemoveConnection(connection types.Connection) {
	m.Called(connection)
}
