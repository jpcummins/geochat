package commands

import (
	"fmt"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

type split struct{}

func (s *split) Execute(args string, user types.User, world types.World) error {
	zoneID := user.ZoneID()
	zone, err := world.Zones().Zone(zoneID)
	if err != nil {
		return err
	}

	zones, err := zone.Split()
	if err != nil {
		return err
	}

	announcement := fmt.Sprintf("Zone '%s' split. New zones: ", zoneID)
	for _, newZone := range zones {
		split := broadcast.Split(zone, newZone)
		if err := newZone.Broadcast(split); err != nil {
			return err
		}

		if newZone.ID() != zone.ID() {
			announcement = announcement + fmt.Sprintf("'%s' (%d users) ", newZone.ID(), newZone.Count())
		}
	}
	return nil
}
