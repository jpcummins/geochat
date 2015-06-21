package events

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type Ping struct{}

func (p *Ping) Type() types.ClientEventType {
	return "ping"
}

func (p *Ping) BeforeBroadcast(event types.ClientEvent) error {
	return errors.New("Pings should not be published")
}
