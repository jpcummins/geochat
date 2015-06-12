package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Leave struct {
	UserID string `json:"user_id"`
}

func (m *Leave) Type() string {
	return "leave"
}

func (l *Leave) BeforePublish(e types.Event) error {
	return nil
}

func (l *Leave) OnReceive(e types.Event) error {
	e.Zone().RemoveUser(l.UserID)
	e.Zone().Broadcast(e)
	return nil
}
