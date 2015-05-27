package chat

type Join struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e *Event) error {
	j.Subscriber.zone.join(j.Subscriber)
	return nil
}
