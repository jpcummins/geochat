package commands

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type merge struct{}

func (m *merge) Execute(args string, user types.User, world types.World) error {

	currentZoneID := user.ZoneID()
	currentZone, err := world.Zones().Zone(currentZoneID)
	if err != nil {
		return err
	}

	if currentZone == nil {
		return errors.New("User is not in a zone")
	}

	if currentZoneID == ":0z" {
		return errors.New("Unable to merge root zone")
	}

	parentZone, err := world.Zones().Zone(currentZone.ParentZoneID())

	if err != nil {
		return err
	}

	if err := parentZone.Merge(); err != nil {
		return err
	}

	return nil
}
