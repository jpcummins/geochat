package chat

type Online struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (o *Online) Type() string {
	return "online"
}
