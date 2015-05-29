package chat

type Online struct {
	User *User `json:"user"`
}

func (o *Online) Type() string {
	return "online"
}

func (o *Online) OnReceive(e *Event) error {
	return nil
}
