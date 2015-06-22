package types

type User interface {
	ID() string
	Name() string
	Location() LatLng

	Zone() Zone
	SetZone(Zone)

	Broadcast(ClientEvent)
	Connect() Connection
	Disconnect(Connection)

	ClientJSON() *ClientUserJSON
	ServerJSON() *ServerUserJSON
	Update(*ServerUserJSON) error
}

type ServerUserJSON struct {
	ID           string `json:"id"`
	CreatedAt    int    `json:"created_at"`
	LastActivity int    `json:"last_activity"`
	Name         string `json:"name"`
	Location     LatLng `json:"location"`
	ZoneID       string `json:"zone_id"`
}

type ClientUserJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
