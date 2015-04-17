package chat

type Leave struct {
	User string `json:"user"`
}

func (m *Leave) Type() string {
	return "leave"
}
