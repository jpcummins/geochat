package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Online struct {
	User *types.User `json:"user"`
}

func (o *Online) Type() string {
	return "online"
}

func (o *Online) BeforePublish(e types.Event) error {
	return nil
}

func (o *Online) OnReceive(e types.Event) error {
	return nil
}
