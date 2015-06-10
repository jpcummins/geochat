package chat

type Split struct {
}

func (s *Split) Type() string {
	return "split"
}

func (s *Split) OnReceive(e *Event, z *Zone) error {
	z.broadcastEvent(e)
	return nil
}
