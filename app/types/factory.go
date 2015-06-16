package types

type Factory interface {
	NewWorld(id string, cache Cache, pubsub PubSub) (World, error)
	NewZone(id string, worldID string, maxUsers int) (Zone, error)
	NewUser(id string, name string, location LatLng) (User, error)
	NewEvent(id string, worldID string, data EventData) (Event, error)
}
