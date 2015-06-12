package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Split struct {
}

func (s *Split) Type() string {
	return "split"
}

func (s *Split) BeforePublish(e types.Event) error {
	return nil
}

func (s *Split) OnReceive(e types.Event) error {
	e.Zone().Broadcast(e)
	return nil
}
