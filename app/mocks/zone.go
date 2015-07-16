package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type Zone struct {
	mock.Mock
}

func (m *Zone) ID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) World() types.World {
	ret := m.Called()

	r0 := ret.Get(0).(types.World)

	return r0
}
func (m *Zone) SouthWest() types.LatLng {
	ret := m.Called()

	r0 := ret.Get(0).(types.LatLng)

	return r0
}
func (m *Zone) NorthEast() types.LatLng {
	ret := m.Called()

	r0 := ret.Get(0).(types.LatLng)

	return r0
}
func (m *Zone) Geohash() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) From() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) To() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) ParentZoneID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) LeftZoneID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) RightZoneID() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Zone) Count() int {
	ret := m.Called()

	r0 := ret.Get(0).(int)

	return r0
}
func (m *Zone) IsOpen() bool {
	ret := m.Called()

	r0 := ret.Get(0).(bool)

	return r0
}
func (m *Zone) SetIsOpen(_a0 bool) {
	m.Called(_a0)
}
func (m *Zone) AddUser(_a0 types.User) {
	m.Called(_a0)
}
func (m *Zone) RemoveUser(_a0 string) {
	m.Called(_a0)
}
func (m *Zone) Broadcast(_a0 types.BroadcastEventData) error {
	ret := m.Called(_a0)
	return ret.Error(0)
}
func (m *Zone) Join(_a0 types.User) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *Zone) Leave(_a0 types.User) error {
	ret := m.Called(_a0)
	return ret.Error(0)
}
func (m *Zone) Message(_a0 types.User, _a1 string) error {
	ret := m.Called(_a0, _a1)
	return ret.Error(0)
}
func (m *Zone) Split() (map[string]types.Zone, error) {
	ret := m.Called()

	var r0 map[string]types.Zone
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(map[string]types.Zone)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Zone) Merge() error {
	ret := m.Called()

	r0 := ret.Error(0)

	return r0
}

func (m *Zone) PubSubJSON() types.PubSubJSON {
	ret := m.Called()
	return ret.Get(0).(types.PubSubJSON)
}

func (m *Zone) Update(json types.PubSubJSON) error {
	ret := m.Called(json)
	return ret.Error(0)
}

func (m *Zone) BroadcastJSON() interface{} {
	ret := m.Called()
	return ret.Get(0)
}
func (m *Zone) UserIDs() []string {
	ret := m.Called()

	var r0 []string
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]string)
	}

	return r0
}
