package types

type Chat interface {
	PubSub() PubSub
	Cache() Cache
	World(id string) (World, error)
	SetWorld(World) error
	GetOrCreateWorld(id string) (World, error)
}
