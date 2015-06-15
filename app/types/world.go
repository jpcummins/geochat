package types

type World interface {
	Factory() Factory         // I don't think this is needed
	MaxUsersForNewZones() int // I don't think this is needed
	GetOrCreateZone(id string) (Zone, error)
	// GetOrCreateZoneForUser(user User) (Zone, error)
}
