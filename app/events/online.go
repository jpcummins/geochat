package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Online struct {
	types.BaseEventData
	User *types.User `json:"user"`
}

func (o *Online) Type() string {
	return "online"
}
