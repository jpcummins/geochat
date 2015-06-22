package types

type DB interface {
	User(userID string, worldID string) (*ServerUserJSON, error)
	SaveUser(*ServerUserJSON) error

	Zone(zoneID string, worldID string) (*ServerZoneJSON, error)
	SaveZone(*ServerZoneJSON) error

	GetWorld(worldID string) (*ServerWorldJSON, error)
	SaveWorld(*ServerWorldJSON) error
}
