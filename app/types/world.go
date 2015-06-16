package types

type World interface {
	GetOrCreateZone(id string) (Zone, error)
	GetOrCreateZoneForUser(user User) (Zone, error)
}
