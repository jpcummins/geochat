package types

type World interface {
	Factory() Factory
	MaxUsersForNewZones() int
	GetOrCreateZone(id string) (Zone, error)
}
