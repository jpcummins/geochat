package chat

type Ping struct {
}

func (p *Ping) Type() string {
	return "ping"
}

func (p *Ping) OnReceive(e *Event) error {
	return nil
}
