package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Offline struct {
	User *types.User `json:"user"`
}

func (o *Offline) Type() string {
	return "offline"
}

func (o *Offline) BeforePublish(e types.Event) error {
	return nil
}

func (o *Offline) OnReceive(e types.Event) error {
	return nil
}
