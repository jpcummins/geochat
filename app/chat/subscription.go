package chat

import (
	"strconv"
	"time"
)

type Subscription struct {
	Id     string      `json:"id"`
	User   *User       `json:"user"`
	Events chan *Event `json:"-"`
	Zone   *Zone       `json:"-"`
}

func CreateSubscription(user *User, zone *Zone) *Subscription {
	return &Subscription{
		Id:     zone.Geohash + user.Id + strconv.Itoa(int(time.Now().Unix())),
		User:   user,
		Zone:   zone,
		Events: make(chan *Event, 10),
	}
}
