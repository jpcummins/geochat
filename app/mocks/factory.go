package mocks

import (
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/mock"
)

type Factory struct {
	mock.Mock
}

func (f *Factory) NewWorld(cache types.Cache, maxUsersForNewZones int) (types.World, error) {
	args := f.Called(cache, maxUsersForNewZones)
	return args.Get(0).(types.World), args.Error(1)
}

func (f *Factory) NewZone(world types.World, id string) (types.Zone, error) {
	args := f.Called(world, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(types.Zone), args.Error(1)
}

func (f *Factory) NewUser(lat float64, long float64, name string, id string) (types.User, error) {
	args := f.Called(lat, long, name, id)
	return args.Get(0).(types.User), args.Error(1)
}
