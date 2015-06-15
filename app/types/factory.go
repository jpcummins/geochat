package types

type Factory interface {
	NewWorld(cache Cache, maxUsersForNewZones int) (World, error)
	NewZone(world World, id string) (Zone, error)
	NewUser(lat float64, long float64, name string, id string) (User, error)
}
