package chat

type Leave struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (m *Leave) Type() string {
	return "leave"
}
