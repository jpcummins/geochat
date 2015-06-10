package chat

type Leave struct {
	UserID string `json:"user_id"`
}

func (m *Leave) Type() string {
	return "leave"
}

func (l *Leave) OnReceive(e *Event, z *Zone) error {
	z.delUser(l.UserID)
	z.broadcastEvent(e)
	return nil
}
