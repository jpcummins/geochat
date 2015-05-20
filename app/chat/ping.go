package chat

type Ping struct {
}

func (j *Ping) Type() string {
	return "ping"
}
