package commands

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type merge struct{}

func (m *merge) Execute(args string, user types.User) error {

	currentZone := user.Zone()

	if currentZone == nil {
		return errors.New("User is not in a zone")
	}

	if currentZone.ID() == ":0z" {
		return errors.New("Unable to merge root zone")
	}

	parentZone, err := currentZone.World().Zones().Zone(currentZone.ParentZoneID())

	if err != nil {
		return err
	}

	return parentZone.Merge()
}
