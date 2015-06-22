package types

type User interface {
	EventJSON

	ID() string
	Name() string
	Location() LatLng

	Zone() Zone
	SetZone(Zone)

	Broadcast(ClientEvent)
	Connect() Connection
	Disconnect(Connection)
}

type ServerUserJSON struct {
	*BaseServerJSON
	CreatedAt    int    `json:"created_at"`
	LastActivity int    `json:"last_activity"`
	Name         string `json:"name"`
	Location     LatLng `json:"location"`
	ZoneID       string `json:"zone_id"`
}

type ClientUserJSON struct {
	*BaseClientJSON
	Name string `json:"name"`
}
