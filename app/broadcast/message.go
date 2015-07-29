package broadcast

import (
	"github.com/jpcummins/geochat/app/types"
)

const messageType types.BroadcastEventType = "message"

type message struct {
	UserID string `json:"userID"`
	Text   string `json:"text"`
}

func Message(userID string, text string) *message {
	return &message{
		UserID: userID,
		Text:   text,
	}
}

func (e *message) Type() types.BroadcastEventType {
	return messageType
}
