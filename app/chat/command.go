package chat

import (
	"encoding/json"
	"strings"
)

type command struct {
	name    string
	usage   string
	execute func([]string, *Subscription) (string, error)
}

var commands = make(map[string]*command)

func registerCommand(command *command) {
	commands[command.name] = command
}

func ExecuteCommand(command string, subscription *Subscription) (string, error) {
	args := strings.Split(command, " ")
	if len(args) == 0 || commands[args[0]] == nil {
		output, err := json.Marshal(commands)
		return string(output[:]), err
	}
	return commands[args[0]].execute(args[1:], subscription)
}

func resetRedis(args []string, subscription *Subscription) (string, error) {
	c := connection.Get()
	defer c.Close()
	c.Do("FLUSHALL")
	return "", nil
}
