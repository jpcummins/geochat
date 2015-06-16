package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Factory struct{}

func (f *Factory) NewWorld(cache types.Cache) (types.World, error) {
	return newWorld(cache, f, 3)
}

func (f *Factory) NewZone(id string, maxUsers int) (types.Zone, error) {
	return newZone(id, maxUsers)
}

func (f *Factory) NewUser(id string, name string, location types.LatLng) (types.User, error) {
	return newUser(id, name, location), nil
}
