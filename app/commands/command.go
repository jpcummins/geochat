package commands

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

func Execute(command string, args string, user types.User) error {
	zone := user.Zone()
	if zone == nil {
		return errors.New("Can't find zone")
	}

	zone.Message(user, "/"+command+" "+args)

	switch command {
	case "addbot":
		return (&addBot{}).Execute(args, user)
	case "split":
		return (&split{}).Execute(args, user)
	}
	return nil
}

type commandName string

type command interface {
	Execute(args string, user types.User) error
}
