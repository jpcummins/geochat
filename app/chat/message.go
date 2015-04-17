package chat

type Message struct {
	User *User  `json:"user"`
	Text string `json:"text"`
}

func (m *Message) Type() string {
	return "message"
}
