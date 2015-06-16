package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Factory struct {
	mock.Mock
}

func (f *Factory) NewWorld(cache types.Cache) (types.World, error) {
	args := f.Called(cache)
	return args.Get(0).(types.World), args.Error(1)
}

func (f *Factory) NewZone(id string, maxUsers int) (types.Zone, error) {
	args := f.Called(id, maxUsers)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(types.Zone), args.Error(1)
}

func (f *Factory) NewUser(id string, name string, location types.LatLng) (types.User, error) {
	args := f.Called(id, name, location)
	return args.Get(0).(types.User), args.Error(1)
}
