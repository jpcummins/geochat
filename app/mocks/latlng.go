package mocks

import (
	"github.com/stretchr/testify/mock"
)

type LatLng struct {
	mock.Mock
}

func (m *LatLng) Lat() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *LatLng) Lng() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *LatLng) Geohash() string {
	args := m.Called()
	return args.String(0)
}
