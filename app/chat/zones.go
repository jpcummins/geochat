package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Zones struct {
	sync.RWMutex
	db    types.DB
	world types.World
	zones map[string]types.Zone
}

func newZones(world types.World, db types.DB) *Zones {
	return &Zones{
		db:    db,
		world: world,
		zones: make(map[string]types.Zone),
	}
}

func (z *Zones) Zone(id string) (types.Zone, error) {
	if cachedZone := z.FromCache(id); cachedZone != nil {
		return cachedZone, nil
	}

	return z.FromDB(id)
}

func (z *Zones) FromCache(id string) types.Zone {
	z.RLock()
	defer z.RUnlock()
	return z.zones[id]
}

func (z *Zones) FromDB(id string) (types.Zone, error) {
	json, err := z.db.Zone(id, z.world.ID())
	if err != nil {
		return nil, err
	}

	if json == nil {
		return nil, nil
	}

	zone := z.FromCache(id)
	if zone == nil {
		zone, err = newZone(id, z.world, z.world.MaxUsers())
		if err != nil {
			return nil, err
		}
	}

	zone.Update(json)
	z.updateCache(zone)
	return zone, nil
}

func (z *Zones) Save(zone types.Zone) error {
	if err := z.db.SaveZone(zone.ServerJSON()); err != nil {
		return err
	}

	z.updateCache(zone)
	return nil
}

func (z *Zones) updateCache(zone types.Zone) {
	z.Lock()
	defer z.Unlock()
	z.zones[zone.ID()] = zone
}
