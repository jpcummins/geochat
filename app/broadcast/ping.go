package broadcast

import (
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

func (e *ping) BeforeBroadcastToUser(user types.User, event types.BroadcastEvent) (bool, error) {
	return true, nil
}
