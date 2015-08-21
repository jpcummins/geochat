package commands

import (
	"errors"
	"github.com/jpcummins/geochat/app/pubsub"
	"github.com/jpcummins/geochat/app/types"
	"strconv"
	"strings"
)

type maxusers struct{}

func (m *maxusers) Execute(args string, user types.User, world types.World) error {
	splitArgs := strings.Split(strings.TrimSpace(args), " ")

	if len(splitArgs) == 1 && splitArgs[0] == "" {
		announcement, err := pubsub.Announcement("Max users: "+strconv.Itoa(world.MaxUsers()), user.ID())
		if err != nil {
			return err
		}
		return world.Publish(announcement)
	}

	if len(splitArgs) != 1 {
		return errors.New("Invalid number of arguments")
	}

	num, err := strconv.Atoi(splitArgs[0])
	if err != nil {
		return err
	}
	world.SetMaxUsers(num)

	json := world.PubSubJSON().(*types.WorldPubSubJSON)
	world.DB().SaveWorld(json)

	pubSubEvent, err := pubsub.World(json)
	if err != nil {
		return err
	}
	world.Publish(pubSubEvent)

	announcement, err := pubsub.Announcement("Max users: "+strconv.Itoa(world.MaxUsers()), user.ID())
	if err != nil {
		return err
	}
	return world.Publish(announcement)
}
