package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Factory struct{}

func (f *Factory) NewWorld(cache types.Cache, maxUsersForNewZones int) (types.World, error) {
	return newWorld(cache, f, maxUsersForNewZones)
}

func (f *Factory) NewZone(world types.World, id string) (types.Zone, error) {
	return newZone(world, id)
}

func (f *Factory) NewUser(lat float64, long float64, name string, id string) (types.User, error) {
	return newUser(lat, long, name, id), nil
}
