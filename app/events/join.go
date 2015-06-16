package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Join struct {
	types.BaseEventData
	userID string
}

func JoinEventData(userID string) (*Join, error) {
	j := &Join{
		userID: userID,
	}
	return j, nil
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e types.Event) error {
	return nil
}
