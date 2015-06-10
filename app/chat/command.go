package chat

type command struct {
	name    string
	usage   string
	execute func([]string, *User) (string, error)
}

var commands = make(map[string]*command)

func registerCommand(command *command) {
	commands[command.name] = command
}

func resetRedis(args []string, user *User) (string, error) {
	c := Redis.Get()
	defer c.Close()
	c.Do("FLUSHALL")
	return "", nil
}

// func join(args []string, subscription *Subscription) (string, error) {
// 	zone, _ := GetOrCreateAvailableZone(subscription.User.Lat, subscription.User.Long)
// 	subscription.zone = zone
// 	zoneInfo := &ZoneInfo{
// 		ID:          zone.Zonehash,
// 		Boundary:    zone.Boundary,
// 		Archive:     nil,
// 		Subscribers: zone.GetSubscribers(),
// 	}
// 	subscription.Events <- NewEvent(zoneInfo)
// 	return "", nil
// }
