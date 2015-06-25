package types

type User interface {
	PubSubSerializable
	BroadcastSerializable

	ID() string
	Name() string
	Location() LatLng

	Zone() Zone
	SetZone(Zone)

	Broadcast(BroadcastEventData)
	ExecuteCommand(command string, args string) error

	Connect() Connection
	Disconnect(Connection)
}

type UserPubSubJSON struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Location *LatLngJSON `json:"location"`
	ZoneID   string      `json:"zone_id"`
}

func (psUser *UserPubSubJSON) Type() PubSubDataType {
	return "user"
}

type UserBroadcastJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
