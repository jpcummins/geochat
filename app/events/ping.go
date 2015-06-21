package events

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

type Ping struct{}

func (p *Ping) Type() string {
	return "ping"
}

func (p *Ping) BeforePublish(event types.Event) error {
	return errors.New("Pings should not be published")
}

func (p *Ping) OnReceive(event types.Event) error {
	return errors.New("Received an unsupported ping event")
}
