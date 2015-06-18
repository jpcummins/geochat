package chat

import (
	"github.com/jpcummins/geochat/app/types"
)

type Chat struct {
	cache               types.Cache
	pubsub              types.PubSub
	maxUsersForNewZones int
}

func newChat(cache types.Cache, pubsub types.PubSub, maxUsersForNewZones int) (types.Chat, error) {
	return &Chat{cache, pubsub, maxUsersForNewZones}, nil
}

func (c *Chat) PubSub() types.PubSub {
	return c.pubsub
}

func (c *Chat) Cache() types.Cache {
	return c.cache
}

func (c *Chat) World(id string) (types.World, error) {
	return c.cache.World(id)
}

func (c *Chat) SetWorld(world types.World) error {
	return c.cache.SetWorld(world)
}

func (c *Chat) GetOrCreateWorld(id string) (types.World, error) {
	return nil, nil
}
