package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Leave struct {
	types.BaseEventData
	UserID string `json:"user_id"`
}

func (m *Leave) Type() string {
	return "leave"
}

func (l *Leave) OnReceive(e types.Event) error {
	return nil
}
