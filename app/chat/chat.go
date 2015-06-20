package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Chat struct {
	db                  types.DB
	pubsub              types.PubSub
	events              types.EventFactory
	world               types.World
	maxUsersForNewZones int
}

func (c *Chat) DB() types.DB {
	return c.db
}

func (c *Chat) PubSub() types.PubSub {
	return c.pubsub
}

func (c *Chat) Events() types.EventFactory {
	return c.events
}

func (c *Chat) World() types.World {
	return c.world
}
