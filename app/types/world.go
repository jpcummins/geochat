package types

type World interface {
	GetOrCreateZone(id string) (*World, error)
}
