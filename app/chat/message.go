package chat

type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func (m *Message) Type() string {
	return "message"
}
