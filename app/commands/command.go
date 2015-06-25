package commands

import (
	"github.com/jpcummins/geochat/app/types"
)

func Execute(command string, args string, user types.User) error {
	switch command {
	case "addbot":
		return (&addBot{}).Execute(args, user)
	}

	return nil
}

type commandName string

type command interface {
	Execute(args string, user types.User) error
}
