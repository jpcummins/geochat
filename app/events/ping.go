package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Ping struct {
	types.BaseEventData
}

func (p *Ping) Type() string {
	return "ping"
}
