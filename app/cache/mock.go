package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type MockUserCache struct {
	mock.Mock
}

func (m *MockUserCache) Get(id string) (types.User, error) {
	args := m.Called(id)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockUserCache) Set(user types.User) error {
	args := m.Called(user)
	return args.Error(0)
}
