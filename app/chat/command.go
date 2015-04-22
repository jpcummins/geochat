package chat

import (
	"encoding/json"
	"strings"
)

type command struct {
	name    string                                 `json:"command"`
	usage   string                                 `json:"usage"`
	execute func([]string, string) (string, error) `json:"-"`
}

var commands = make(map[string]*command)

func registerCommand(command *command) {
	commands[command.name] = command
}

func ExecuteCommand(command string, geohash string) (string, error) {
	args := strings.Split(command, " ")
	if len(args) == 0 || commands[args[0]] == nil {
		output, err := json.Marshal(commands)
		return string(output[:]), err
	}
	return commands[args[0]].execute(args[1:], geohash)
}

func resetRedis(args []string, from_geohash string) (string, error) {
	c := pool.Get()
	defer c.Close()
	c.Do("FLUSHALL")
	return "", nil
}
