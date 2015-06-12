package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Offline struct {
	types.BaseEventData
	User *types.User `json:"user"`
}

func (o *Offline) Type() string {
	return "offline"
}
