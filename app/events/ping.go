package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Ping struct {
}

func (p *Ping) Type() string {
	return "ping"
}

func (p *Ping) OnReceive(e types.Event) error {
	return nil
}
