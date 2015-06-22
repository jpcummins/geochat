package types

type ServerJoinJSON struct {
	ZoneID string `json:"zone_id"`
	UserID string `json:"user_id"`
}

type ClientJoinJSON struct {
	ZoneID string         `json:"zone_id"`
	User   ClientUserJSON `json:"user"`
}
