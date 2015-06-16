package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Split struct {
	types.BaseEventData
}

func (s *Split) Type() string {
	return "split"
}

func (s *Split) OnReceive(e types.Event) error {
	return nil
}
