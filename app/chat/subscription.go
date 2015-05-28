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
	id           string
	zone         *Zone
	user         *User
	createdAt    int
	lastActivity int
	isOnline     bool
	Events       chan *Event
}

type subscriptionJSON struct {
	ID           string `json:"id"`
	User         *User  `json:"user"`
	CreatedAt    int    `json:"created_at"`
	LastActivity int    `json:"last_activity"`
	IsOnline     bool   `json:"is_online"`
	Zone         string `json:"zone"`
}

// NewLocalSubscription is a factory method for creating new local subscriptions.
func NewLocalSubscription(user *User) (*Subscription, error) {
	subscription := &Subscription{
		id:           strconv.Itoa(rand.Intn(1000)) + strconv.Itoa(int(time.Now().Unix())),
		user:         user,
		createdAt:    int(time.Now().Unix()),
		lastActivity: int(time.Now().Unix()),
	}

	zone, err := getOrCreateAvailableZone(subscription.user.Lat, subscription.user.Long)
	if err != nil {
		return nil, err
	}

	// This indirectly adds the subscription to the zone's subscriber list. The
	// event is announced, then picked up by onJoinEvent.
	zone.Publish(NewEvent(&Join{Subscriber: subscription}))
	subscription.zone = zone

	// Save the subscription to Redis
	c := connection.Get()
	defer c.Close()
	subscriptionJSON, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}
	c.Do("LPUSH", "subscribers_"+zone.id, subscriptionJSON)

	Subscribers.Set(subscription)
	return subscription, err
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
	s.lastActivity = js.LastActivity
	s.isOnline = js.IsOnline

	zone, err := GetOrCreateZone(js.Zone)
	if err != nil {
		return err
	}
	s.zone = zone

	if !s.zone.isInitialized() {
		s.zone.initialize()
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
		ID:           s.id,
		User:         s.user,
		CreatedAt:    s.createdAt,
		LastActivity: s.lastActivity,
		IsOnline:     s.isOnline,
		Zone:         zoneID,
	}

	return json.Marshal(subscriptionJSON)
}

// GetZone returns the zone associated to the subscription
func (s *Subscription) GetZone() *Zone {
	return s.zone
}

// GetUser returns the user associated to the subscription
func (s *Subscription) GetUser() *User {
	return s.user
}

// GetID returns the current subscription id
func (s *Subscription) GetID() string {
	return s.id
}

// IsLocal returns true if the subscription is connected to this server instance
func (s *Subscription) IsLocal() bool {
	return s.Events != nil
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

// UpdateLastActiveTime sets the last active time for the subscriber
func (s *Subscription) UpdateLastActiveTime() {
	s.lastActivity = int(time.Now().Unix())
	Subscribers.Set(s)
}
