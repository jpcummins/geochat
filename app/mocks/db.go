package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type DB struct {
	mock.Mock
}

func (m *DB) Publish(event types.Event) error {
	ret := m.Called(event)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) Subscribe(zone types.Zone) <-chan types.Event {
	ret := m.Called(zone)

	var r0 <-chan types.Event
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(<-chan types.Event)
	}

	return r0
}
func (m *DB) GetUser(id string) (types.User, error) {
	ret := m.Called(id)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetUser(user types.User) error {
	ret := m.Called(user)

	r0 := ret.Error(0)

	return r0
}
func (m *DB) GetZone(id string) (types.Zone, error) {
	ret := m.Called(id)

	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DB) SetZone(zone types.Zone) error {
	ret := m.Called(zone)

	r0 := ret.Error(0)

	return r0
}
