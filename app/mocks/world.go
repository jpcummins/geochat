package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

import "time"

type World struct {
	mock.Mock
}

func (_m *World) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *World) MaxUsers() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
func (_m *World) SetMaxUsers(_a0 int) {
	_m.Called(_a0)
}
func (_m *World) MinUsers() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
func (_m *World) SetMinUsers(_a0 int) {
	_m.Called(_a0)
}
func (_m *World) SplitDelay() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}
func (_m *World) DB() types.DB {
	ret := _m.Called()

	var r0 types.DB
	if rf, ok := ret.Get(0).(func() types.DB); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.DB)
	}

	return r0
}
func (_m *World) Zone() types.Zone {
	ret := _m.Called()

	var r0 types.Zone
	if rf, ok := ret.Get(0).(func() types.Zone); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.Zone)
	}

	return r0
}
func (_m *World) Zones() types.Zones {
	ret := _m.Called()

	var r0 types.Zones
	if rf, ok := ret.Get(0).(func() types.Zones); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.Zones)
	}

	return r0
}
func (_m *World) GetOrCreateZone(_a0 string) (types.Zone, error) {
	ret := _m.Called(_a0)

	var r0 types.Zone
	if rf, ok := ret.Get(0).(func(string) types.Zone); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.Zone)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *World) FindOpenZone(_a0 types.Zone, _a1 types.User) (types.Zone, error) {
	ret := _m.Called(_a0, _a1)

	var r0 types.Zone
	if rf, ok := ret.Get(0).(func(types.Zone, types.User) types.Zone); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(types.Zone)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.Zone, types.User) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *World) Join(_a0 types.User) (types.Zone, error) {
	ret := _m.Called(_a0)

	var r0 types.Zone
	if rf, ok := ret.Get(0).(func(types.User) types.Zone); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.Zone)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *World) Users() types.Users {
	ret := _m.Called()

	var r0 types.Users
	if rf, ok := ret.Get(0).(func() types.Users); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.Users)
	}

	return r0
}
func (_m *World) NewUser(_a0 string) (types.User, error) {
	ret := _m.Called(_a0)

	var r0 types.User
	if rf, ok := ret.Get(0).(func(string) types.User); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *World) Publish(_a0 types.PubSubEventData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.PubSubEventData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *World) PubSubJSON() types.PubSubJSON {
	ret := _m.Called()

	var r0 types.PubSubJSON
	if rf, ok := ret.Get(0).(func() types.PubSubJSON); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.PubSubJSON)
	}

	return r0
}
func (_m *World) Update(_a0 types.PubSubJSON) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.PubSubJSON) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
