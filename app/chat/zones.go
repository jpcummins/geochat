package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Zones struct {
	sync.RWMutex
	zones map[string]types.Zone
}

func newZones() *Zones {
	return &Zones{
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
	found, err := z.dependencies.DB().GetZone(id, w, zone)
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
	if err := z.dependencies.DB().SetZone(zone, w); err != nil {
		return err
	}

	z.localSetZone(zone)
	return nil
}

func (z *Zones) localZone(id string) (types.Zone, bool) {
	z.zoneMutex.RLock()
	defer z.zoneMutex.RUnlock()
	zone, found := z.zones[id]
	return zone, found
}

func (z *Zones) localSetZone(zone types.Zone) {
	z.zoneMutex.Lock()
	defer z.zoneMutex.Unlock()
	z.zones[zone.ID()] = zone
}
