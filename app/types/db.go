package types

type DB interface {
	User(userID string, worldID string) (*UserPubSubJSON, error)
	SaveUser(json *UserPubSubJSON, worldID string) error

	Zone(zoneID string, worldID string) (*ZonePubSubJSON, error)
	SaveZone(json *ZonePubSubJSON, worldID string) error

	SaveUsersAndZones(users []*UserPubSubJSON, zones []*ZonePubSubJSON, worldID string) error

	World(worldID string) (*WorldPubSubJSON, error)
	SaveWorld(json *WorldPubSubJSON) error
}
