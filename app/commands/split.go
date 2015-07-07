package commands

import (
	"github.com/jpcummins/geochat/app/types"
)

type split struct{}

func (s *split) Execute(args string, user types.User) error {
	event, err := user.Zone().Split()
	if err != nil {
		return err
	}

	return user.Zone().Broadcast(event)
}
