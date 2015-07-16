package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
	log "gopkg.in/inconshreveable/log15.v2"
	"sync"
)

type Zones struct {
	sync.RWMutex
	db     types.DB
	world  *World
	zones  map[string]types.Zone
	logger log.Logger
}

func newZones(world *World, db types.DB, logger log.Logger) *Zones {
	return &Zones{
		db:     db,
		world:  world,
		zones:  make(map[string]types.Zone),
		logger: logger,
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
		z.logger.Error("Error finding zone from db", "zone", id, "error", err.Error())
		return nil, err
	}

	if json == nil {
		return nil, nil
	}

	zone := z.FromCache(id)
	if zone == nil {
		if zone, err = newZone(id, z.world, z.logger); err != nil {
			z.logger.Error("Error creating zone", "zone", id, "error", err.Error())
			return nil, err
		}
	}
	if err := zone.Update(json); err != nil {
		z.logger.Error("Error updating zone", "zone", id, "error", err.Error())
		return nil, err
	}

	return zone, nil
}

func (z *Zones) Save(zone types.Zone) error {
	json, ok := zone.PubSubJSON().(*types.ZonePubSubJSON)
	if !ok {
		return errors.New("Unable to serialize ZonePubSubJSON")
	}

	if err := z.db.SaveZone(json, z.world.ID()); err != nil {
		return err
	}

	z.UpdateCache(zone)
	return nil
}

func (z *Zones) UpdateCache(zone types.Zone) {
	z.Lock()
	defer z.Unlock()
	z.zones[zone.ID()] = zone
}
