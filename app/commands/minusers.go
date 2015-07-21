package commands

import (
	"errors"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/pubsub"
	"github.com/jpcummins/geochat/app/types"
	"strconv"
	"strings"
)

type minusers struct{}

func (m *minusers) Execute(args string, user types.User) error {
	splitArgs := strings.Split(args, " ")

	if len(splitArgs) != 1 {
		return errors.New("Invalid number of arguments")
	}

	num, err := strconv.Atoi(splitArgs[0])
	if err != nil {
		return err
	}
	user.Zone().World().SetMinUsers(num)

	json := user.Zone().World().PubSubJSON().(*types.WorldPubSubJSON)
	user.Zone().World().DB().SaveWorld(json)
	_, err = pubsub.World(json)
	user.Broadcast(broadcast.Message(user.ID(), "World Updated"))
	return err
}
