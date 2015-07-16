package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Zones struct {
	mock.Mock
}

func (m *Zones) Zone(_a0 string) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Zones) FromCache(_a0 string) types.Zone {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)

	return r0
}
func (m *Zones) UpdateCache(_a0 types.Zone) {
	m.Called(_a0)
}
func (m *Zones) FromDB(_a0 string) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Zones) Save(_a0 types.Zone) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
