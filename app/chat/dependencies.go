package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Dependencies struct {
	db                  types.DB
	pubsub              types.PubSub
	events              types.EventFactory
	maxUsersForNewZones int
}

func (c *Dependencies) DB() types.DB {
	return c.db
}

func (c *Dependencies) PubSub() types.PubSub {
	return c.pubsub
}

func (c *Dependencies) Events() types.EventFactory {
	return c.events
}
