package commands

import (
	"errors"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

func Execute(command string, args string, user types.User, world types.World) error {
	if user.ZoneID() == "" {
		return errors.New("Can't find zone")
	}

	user.Broadcast(broadcast.Message(user.ID(), "/"+command+" "+args))

	switch command {
	case "addbot":
		return (&addBot{}).Execute(args, user, world)
	case "split":
		return (&split{}).Execute(args, user, world)
	case "merge":
		return (&merge{}).Execute(args, user, world)
	case "minusers":
		return (&minusers{}).Execute(args, user, world)
	case "maxusers":
		return (&maxusers{}).Execute(args, user, world)
	}
	return nil
}

type commandName string

type command interface {
	Execute(args string, user types.User, world types.World) error
}
