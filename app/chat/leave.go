package chat

type Leave struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (m *Leave) Type() string {
	return "leave"
}

func (l *Leave) OnReceive(e *Event) error {
	return nil
}
