package types

type Factory interface {
	NewWorld(cache Cache) (World, error)
	NewZone(id string, maxUsers int) (Zone, error)
	NewUser(id string, name string, location LatLng) (User, error)
}
