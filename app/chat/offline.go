package chat

type Offline struct {
	User *User `json:"user"`
}

func (o *Offline) Type() string {
	return "offline"
}

func (o *Offline) OnReceive(e *Event) error {
	return nil
}
