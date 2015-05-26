package chat

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

type Subscription struct {
	Id        string      `json:"id"`
	User      *User       `json:"user"`
	CreatedAt int         `json:"created_at"`
	IsOnline  bool        `json:"is_online"`
	Zonehash  string      `json:"zonehash"`
	IsLocal   bool        `json:"-"`
	Events    chan *Event `json:"-"`
	zone      *Zone       `json:"-"`
}

type jsonSubscription struct {
	Id        string `json:"id"`
	User      *User  `json:"user"`
	CreatedAt int    `json:"created_at"`
	IsOnline  bool   `json:"is_online"`
	Zonehash  string `json:"zonehash"`
}

func NewSubscription(user *User, zone *Zone) *Subscription {
	subscription := &Subscription{
		Id:   strconv.Itoa(rand.Intn(1000)) + strconv.Itoa(int(time.Now().Unix())),
		User: user,
		zone: zone,
	}
	return subscription
}

func (s *Subscription) UnmarshalJSON(b []byte) error {
	var js jsonSubscription
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}
	s.Id = js.Id
	s.User = js.User
	s.CreatedAt = js.CreatedAt
	s.IsOnline = js.IsOnline
	s.Zonehash = js.Zonehash
	s.IsLocal = false

	zone, err := GetOrCreateZone(s.Zonehash)
	s.zone = zone
	return err
}

func (s *Subscription) Activate() {
	s.IsLocal = true
	s.zone.Publish(NewEvent(&Online{s}))
	subscribers.publishOnline <- s
}

func (s *Subscription) Deactivate() {
	s.IsLocal = false

	if s.Events != nil {
		close(s.Events)
		s.Events = nil
	}

	if s.IsOnline {
		s.IsOnline = false
		s.zone.Publish(NewEvent(&Offline{s}))
		subscribers.publishOffline <- s
	}
}

func (s *Subscription) Broadcast(text string) *Event {
	m := &Message{User: s.User, Text: text}
	e := NewEvent(m)
	s.zone.Publish(e)
	return e
}

func (s *Subscription) SetZone(zone *Zone) err {
	return nil
}
