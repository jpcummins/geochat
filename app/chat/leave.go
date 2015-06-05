package chat

type Leave struct {
	UserID string `json:"user_id"`
	ZoneID string `json:"zone_id"`
}

func (m *Leave) Type() string {
	return "leave"
}

func (l *Leave) OnReceive(e *Event) error {
	zone, err := GetOrCreateZone(l.ZoneID)

	if err != nil {
		zone.delUser(l.UserID)
		zone.broadcastEvent(e)
	}

	return err
}
