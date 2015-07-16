package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const leaveType types.BroadcastEventType = "leave"

type leave struct {
	UserID interface{} `json:"user_id"`
}

func Leave(user types.User) *leave {
	return &leave{
		UserID: user.ID(),
	}
}

func (e *leave) Type() types.BroadcastEventType {
	return leaveType
}
