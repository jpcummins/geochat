package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Users struct {
	mock.Mock
}

func (m *Users) User(_a0 string) (types.User, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Users) FromCache(_a0 string) types.User {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.User)

	return r0
}
func (m *Users) UpdateCache(_a0 types.User) {
	m.Called(_a0)
}
func (m *Users) FromDB(_a0 string) (types.User, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.User)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Users) Save(_a0 types.User) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
