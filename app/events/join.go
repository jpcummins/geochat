package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Join struct {
	userID string
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) BeforePublish(e types.Event) error {
	return nil
}

func (j *Join) OnReceive(e types.Event) error {
	// e.Zone().AddUser(j.user)
	e.Zone().Broadcast(e)
	return nil
}
