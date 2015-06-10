package chat

type Split struct {
}

func (s *Split) Type() string {
	return "split"
}

func (s *Split) OnReceive(e *Event) error {
	return nil
}
