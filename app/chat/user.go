package chat

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type User struct {
	id           string
	zone         *Zone
	createdAt    int
	lastActivity int
	isOnline     bool
	name         string
	lat          float64
	long         float64
	Events       chan *Event
}

type userJSON struct {
	ID           string `json:"id"`
	CreatedAt    int    `json:"created_at"`
	LastActivity int    `json:"last_activity"`
	IsOnline     bool   `json:"is_online"`
	Zone         string `json:"zone"`
	Name         string `json:"name"`
}

func NewLocalUser(lat float64, long float64, name string) (*User, error) {
	user := &User{
		id:           strconv.Itoa(rand.Intn(1000)) + strconv.Itoa(int(time.Now().Unix())),
		createdAt:    int(time.Now().Unix()),
		lastActivity: int(time.Now().Unix()),
		name:         name,
	}

	zone, err := getOrCreateAvailableZone(lat, long)
	if err != nil {
		return nil, err
	}

	zone.AddUser(user)
	return user, err
}

func (u *User) UnmarshalJSON(b []byte) error {
	var js userJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	if _, found := UserCache.cacheGet(js.ID); found {
		panic(errors.New("Attempted to unmarshal a known user."))
	}

	u.id = js.ID
	u.createdAt = js.CreatedAt
	u.lastActivity = js.LastActivity
	u.isOnline = js.IsOnline
	u.name = js.Name

	zone, err := GetOrCreateZone(js.Zone)
	if err != nil {
		return err
	}
	u.zone = zone

	if !u.zone.isInitialized() {
		u.zone.initialize()
	}

	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {

	var zoneID string
	if u.zone == nil {
		zoneID = ""
	} else {
		zoneID = u.zone.GetID()
	}

	userJSON := &userJSON{
		ID:           u.id,
		CreatedAt:    u.createdAt,
		LastActivity: u.lastActivity,
		IsOnline:     u.isOnline,
		Zone:         zoneID,
		Name:         u.name,
	}

	return json.Marshal(userJSON)
}

// GetZone returns the zone associated to the subscription
func (u *User) GetZone() *Zone {
	return u.zone
}

// GetID returns the current subscription id
func (u *User) GetID() string {
	return u.id
}

// IsConnected returns true if the subscription is connected to this server instance
func (u *User) IsConnected() bool {
	return u.Events != nil
}

func (u *User) Connect() {
	u.Events = make(chan *Event, 10)
	u.isOnline = true
	UserCache.Set(u)
}

func (u *User) Disconnect() {
	close(u.Events)
	u.Events = nil
	u.isOnline = false
	UserCache.Set(u)
}

// ExecuteCommand allows certain subscribers to issue administrative commands.
func (u *User) ExecuteCommand(command string) (string, error) {
	args := strings.Split(command, " ")
	if len(args) == 0 || commands[args[0]] == nil {
		output, err := json.Marshal(commands)
		return string(output[:]), err
	}
	return commands[args[0]].execute(args[1:], u)
}

// UpdateLastActiveTime sets the last active time for the subscriber
func (u *User) UpdateLastActiveTime() {
	u.lastActivity = int(time.Now().Unix())
	UserCache.Set(u)
}
