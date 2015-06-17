package types

type Chat interface {
	World(id string) (World, error)
	SetWorld(World) error
	GetOrCreateWorld(id string) (World, error)
}
