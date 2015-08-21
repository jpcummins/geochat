package pubsub

import (
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/types"
)

const announcementType types.PubSubEventType = "announcement"

type announcement struct {
	Message  string `json:"message"`
	ToUserID string `json:"user_id"`
}

func Announcement(message string, toUserID string) (*announcement, error) {
	a := &announcement{
		Message:  message,
		ToUserID: toUserID,
	}
	return a, nil
}

func (a *announcement) Type() types.PubSubEventType {
	return announcementType
}

func (a *announcement) BeforePublish(e types.PubSubEvent) error {
	return nil
}

func (a *announcement) OnReceive(e types.PubSubEvent) error {
	user, err := e.World().Users().User(a.ToUserID)
	if err != nil {
		return err
	}

	user.Broadcast(broadcast.Announcement(a.Message))
	return nil
}
