package mocks

import "github.com/jpcummins/geochat/app/types"
import "github.com/stretchr/testify/mock"

type User struct {
	mock.Mock
}

func (_m *User) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetName(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) FirstName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetFirstName(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) LastName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetLastName(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) Timezone() float64 {
	ret := _m.Called()

	var r0 float64
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	return r0
}
func (_m *User) SetTimezone(_a0 float64) {
	_m.Called(_a0)
}
func (_m *User) Email() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetEmail(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) Location() types.LatLng {
	ret := _m.Called()

	var r0 types.LatLng
	if rf, ok := ret.Get(0).(func() types.LatLng); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.LatLng)
	}

	return r0
}
func (_m *User) SetLocation(_a0 float64, _a1 float64) {
	_m.Called(_a0, _a1)
}
func (_m *User) ZoneID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetZoneID(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) FBID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetFBID(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) FBAccessToken() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetFBAccessToken(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) FBPictureURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *User) SetFBPictureURL(_a0 string) {
	_m.Called(_a0)
}
func (_m *User) Broadcast(_a0 types.BroadcastEventData) {
	_m.Called(_a0)
}
func (_m *User) ExecuteCommand(command string, args string) error {
	ret := _m.Called(command, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(command, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *User) Connect() types.Connection {
	ret := _m.Called()

	var r0 types.Connection
	if rf, ok := ret.Get(0).(func() types.Connection); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.Connection)
	}

	return r0
}
func (_m *User) Disconnect(_a0 types.Connection) {
	_m.Called(_a0)
}

func (_m *User) PubSubJSON() types.PubSubJSON {
	ret := _m.Called()

	var r0 types.PubSubJSON
	if rf, ok := ret.Get(0).(func() types.PubSubJSON); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.PubSubJSON)
	}

	return r0
}
func (_m *User) Update(_a0 types.PubSubJSON) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.PubSubJSON) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *User) BroadcastJSON() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}
