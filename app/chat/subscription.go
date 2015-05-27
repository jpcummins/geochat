package chat

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Subscription represesnts a user's connection. The connection may or may not
// be local to this server instance. "events" is nil for remote connections,
// connections where the user's websocket is handled by another server instance.
type Subscription struct {
	id        string
	zone      *Zone
	user      *User
	createdAt int
	isOnline  bool
	Events    chan *Event
}

type subscriptionJSON struct {
	ID        string `json:"id"`
	User      *User  `json:"user"`
	CreatedAt int    `json:"created_at"`
	IsOnline  bool   `json:"is_online"`
	Zone      string `json:"zone"`
}

// NewSubscription is a factory method for creating new subscriptions.
func NewSubscription(user *User) *Subscription {
	subscription := &Subscription{
		id:        strconv.Itoa(rand.Intn(1000)) + strconv.Itoa(int(time.Now().Unix())),
		user:      user,
		createdAt: int(time.Now().Unix()),
	}
	subscribers.publishSubscribe <- subscription
	return subscription
}

// UnmarshalJSON overrides Go's default JSON unmarshaling method so that this
// object can be marshaled/unmarshaled as the type `subscriptionJSON`.
func (s *Subscription) UnmarshalJSON(b []byte) error {
	var js subscriptionJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	s.id = js.ID
	s.user = js.User
	s.createdAt = js.CreatedAt
	s.isOnline = js.IsOnline

	if js.Zone != "" {
		zone, err := GetOrCreateZone(js.Zone)
		if err != nil {
			return err
		}
		s.zone = zone
	}

	return nil
}

// MarshalJSON overrides Go's default JSON marshaling method so that this object
// can be marshaled/unmarshaled as the type `subscriptionJSON`
func (s *Subscription) MarshalJSON() ([]byte, error) {

	var zoneID string
	if s.zone == nil {
		zoneID = ""
	} else {
		zoneID = s.zone.GetID()
	}

	subscriptionJSON := &subscriptionJSON{
		ID:        s.id,
		User:      s.user,
		CreatedAt: s.createdAt,
		IsOnline:  s.isOnline,
		Zone:      zoneID,
	}

	return json.Marshal(subscriptionJSON)
}

// SetOnline sets the subscription to an online state and communicates the
// state change to other server nodes.
func (s *Subscription) SetOnline() {
	s.isOnline = true
	s.Events = make(chan *Event, 10)
	subscribers.publishOnline <- s
}

// SetOffline sets the subscription to an offline state and communicates the
// state change to other server nodes.
func (s *Subscription) SetOffline() {
	close(s.Events)
	s.Events = nil
	s.isOnline = false
	subscribers.publishOffline <- s
}

// GetZone returns the user's current zone
func (s *Subscription) GetZone() *Zone {
	return s.zone
}

// GetUser returns the user associated to the subscription
func (s *Subscription) GetUser() *User {
	return s.user
}

// SetZone is used to change the subscriber's zone and communicate it to other
// nodes on the network.
func (s *Subscription) SetZone(zone *Zone) {
	s.zone = zone
}

// GetID returns the current subscription id
func (s *Subscription) GetID() string {
	return s.id
}

// ExecuteCommand allows certain subscribers to issue administrative commands.
func (s *Subscription) ExecuteCommand(command string) (string, error) {
	args := strings.Split(command, " ")
	if len(args) == 0 || commands[args[0]] == nil {
		output, err := json.Marshal(commands)
		return string(output[:]), err
	}
	return commands[args[0]].execute(args[1:], s)
}
