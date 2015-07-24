package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const announcementType types.BroadcastEventType = "announcement"

type announcement struct {
	Text string `json:"text"`
}

func Announcement(text string) *announcement {
	return &announcement{
		Text: text,
	}
}

func (a *announcement) Type() types.BroadcastEventType {
	return announcementType
}
