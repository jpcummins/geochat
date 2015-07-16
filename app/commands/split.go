package commands

import (
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

type split struct{}

func (s *split) Execute(args string, user types.User) error {
	prvZone := user.Zone()

	zones, err := user.Zone().Split()
	if err != nil {
		return err
	}

	for _, zone := range zones {
		split := broadcast.Split(prvZone, zone)
		if err := zone.Broadcast(split); err != nil {
			return err
		}
	}

	return nil
}
