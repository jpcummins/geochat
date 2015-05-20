package chat

import (
	"encoding/json"
)

type Subscription struct {
	Id        string      `json:"id"`
	User      *User       `json:"user"`
	CreatedAt int         `json:"created_at"`
	IsOnline  bool        `json:"is_online"`
	Events    chan *Event `json:"-"`
	Zonehash  string      `json:"zonehash"`
	zone      *Zone       `json:"-"`
}

type jsonSubscription struct {
	Id        string `json:"id"`
	User      *User  `json:"user"`
	CreatedAt int    `json:"created_at"`
	IsOnline  bool   `json:"is_online"`
	Zonehash  string `json:"zonehash"`
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

	zone, err := GetOrCreateZone(s.Zonehash)
	s.zone = zone
	return err
}

func (s *Subscription) Activate() {
	if !s.IsOnline {
		s.Events = make(chan *Event, 10)
		s.IsOnline = true
		s.zone.Publish(NewEvent(&Online{s}))
	}
}

func (s *Subscription) Deactivate() {
	if s.IsOnline {
		s.IsOnline = false
		close(s.Events)
		s.Events = nil
		s.zone.Publish(NewEvent(&Offline{s}))
	}
}

func (s *Subscription) Broadcast(text string) *Event {
	m := &Message{User: s.User, Text: text}
	e := NewEvent(m)
	s.zone.Publish(e)
	return e
}
