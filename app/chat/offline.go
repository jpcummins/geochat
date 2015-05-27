package chat

type Offline struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (o *Offline) Type() string {
	return "offline"
}

func (o *Offline) OnReceive(e *Event) error {
	return nil
}
