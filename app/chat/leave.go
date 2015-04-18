package chat

type Leave struct {
	User *User `json:"user"`
}

func (m *Leave) Type() string {
	return "leave"
}
