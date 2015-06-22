package types

type DB interface {
	User(userID string, worldID string) (*ServerUserJSON, error)
	SaveUser(json ServerJSON) error

	Zone(zoneID string, worldID string) (*ServerZoneJSON, error)
	SaveZone(json ServerJSON) error

	World(worldID string) (*ServerWorldJSON, error)
	SaveWorld(ServerJSON) error
}
