package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Join struct {
	types.BaseEventData
	userID string
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e types.Event) error {
	return nil
}
