package broadcast

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
)

const joinType types.BroadcastEventType = "join"

type join struct {
	UserJSON *types.UserBroadcastJSON `json:"user"`
	user     types.User
}

func Join(user types.User) (*join, error) {
	userJSON, ok := user.BroadcastJSON().(*types.UserBroadcastJSON)
	if !ok {
		return nil, errors.New("Unable to serialize UserBroadcastJSON")
	}

	return &join{
		UserJSON: userJSON,
		user:     user,
	}, nil
}

func (e *join) Type() types.BroadcastEventType {
	return joinType
}

func (e *join) BeforeBroadcast(event types.BroadcastEvent) error {
	return nil
}
