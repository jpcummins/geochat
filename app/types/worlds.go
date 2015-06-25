package types

type Worlds interface {
	World(string) (World, error)
	FromCache(string) World
	FromDB(string) (World, error)
	Save(World) error
}
