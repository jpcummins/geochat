package broadcast

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

const pingType types.BroadcastEventType = "ping"

type ping struct{}

func Ping() *ping {
	return &ping{}
}

func (e *ping) Type() types.BroadcastEventType {
	return pingType
}

func (e *ping) Data() *ping {
	return nil
}

func (p *ping) BeforeBroadcast(event types.BroadcastEvent) error {
	return errors.New("Pings should not be published")
}
