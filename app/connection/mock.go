package connection

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Publish(event types.Event) error {
	return nil
}

func (m *Mock) Subscribe(zone types.Zone) <-chan types.Event {
	return nil
}

func (m *Mock) GetUser(id string) (types.User, error) {
	return nil, nil
}

func (m *Mock) SetUser(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *Mock) GetZone(id string) (types.Zone, error) {
	return nil, nil
}

func (m *Mock) SetZone(zone types.Zone) error {
	args := m.Called(zone)
	return args.Error(0)
}
