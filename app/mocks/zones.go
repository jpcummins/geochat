package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Zones struct {
	mock.Mock
}

func (m *Zones) Zone(_a0 string) (types.Zone, error) {
	ret := m.Called(_a0)

	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Zones) UpdateZone(_a0 string) (types.Zone, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(types.Zone)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Zones) SetZone(_a0 types.Zone) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
