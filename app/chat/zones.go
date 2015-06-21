package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Zones struct {
	sync.RWMutex
	world *World
	zones map[string]types.Zone
}

func newZones(world *World) *Zones {
	return &Zones{
		world: world,
		zones: make(map[string]types.Zone),
	}
}

func (z *Zones) Zone(id string) (types.Zone, error) {
	zone, found := z.localZone(id)
	if found {
		return zone, nil
	}
	return z.UpdateZone(id)
}

func (z *Zones) UpdateZone(id string) (types.Zone, error) {
	zone := &Zone{}
	found, err := z.world.db.GetZone(id, z.world, zone)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	z.localSetZone(zone)
	return zone, nil
}

func (z *Zones) SetZone(zone types.Zone) error {
	if err := z.world.db.SetZone(zone, z.world); err != nil {
		return err
	}

	z.localSetZone(zone)
	return nil
}

func (z *Zones) localZone(id string) (types.Zone, bool) {
	z.RLock()
	defer z.RUnlock()
	zone, found := z.zones[id]
	return zone, found
}

func (z *Zones) localSetZone(zone types.Zone) {
	z.Lock()
	defer z.Unlock()
	z.zones[zone.ID()] = zone
}
