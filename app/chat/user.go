package chat

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User struct {
	sync.RWMutex
	id           string
	zone         *Zone
	createdAt    int
	lastActivity int
	isOnline     bool
	name         string
	lat          float64
	long         float64
	connections  []*Connection
}

type userJSON struct {
	ID           string `json:"id"`
	CreatedAt    int    `json:"created_at"`
	LastActivity int    `json:"last_activity"`
	IsOnline     bool   `json:"is_online"`
	ZoneID       string `json:"zone_id"`
	Name         string `json:"name"`
}

var r = rand.New(rand.NewSource(342324))

func NewLocalUser(lat float64, long float64, name string) (*User, error) {
	user := &User{
		id:           name + strconv.Itoa(r.Intn(1000000)),
		createdAt:    int(time.Now().Unix()),
		lastActivity: int(time.Now().Unix()),
		name:         name,
		connections:  make([]*Connection, 0),
	}

	zone, err := getOrCreateAvailableZone(lat, long)
	if err != nil {
		return nil, err
	}

	user.Join(zone)
	return user, err
}

func (u *User) UnmarshalJSON(b []byte) error {
	var js userJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	if _, found := UserCache.cacheGet(js.ID); found {
		panic(errors.New("Attempted to unmarshal a known user: " + js.ID))
	}

	u.id = js.ID
	u.createdAt = js.CreatedAt
	u.lastActivity = js.LastActivity
	u.isOnline = js.IsOnline
	u.name = js.Name

	zone, err := GetOrCreateZone(js.ZoneID)
	if err != nil {
		return err
	}

	u.zone = zone
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	userJSON := &userJSON{
		ID:           u.id,
		CreatedAt:    u.createdAt,
		LastActivity: u.lastActivity,
		IsOnline:     u.isOnline,
		ZoneID:       u.zone.GetID(),
		Name:         u.name,
	}

	return json.Marshal(userJSON)
}

// GetZone returns the zone associated to the subscription
func (u *User) GetZone() *Zone {
	return u.zone
}

func (u *User) SetOnline() {
	u.GetZone().Publish(NewEvent(&Online{User: u}))
}

func (u *User) SetOffline() {
	u.GetZone().Publish(NewEvent(&Offline{User: u}))
}

func (u *User) Join(z *Zone) {
	u.isOnline = true
	u.zone = z
	u.zone.join <- u
	u.zone.Publish(NewEvent(&Join{User: u}))
}

func (u *User) Leave() {
	u.isOnline = false
	u.zone.leave <- u
	u.zone.Publish(NewEvent(&Leave{UserID: u.GetID(), ZoneID: u.zone.id}))

	u.zone = nil
	UserCache.Set(u)
}

// GetID returns the current subscription id
func (u *User) GetID() string {
	return u.id
}

func (u *User) Connect() *Connection {
	c := newConnection(u)
	u.Lock()
	u.connections = append(u.connections, c)
	u.Unlock()
	return c
}

func (u *User) Disconnect(c *Connection) {
	close(c.Events)
	c.Events = nil

	u.Lock()
	for i, connection := range u.connections {
		if connection == c {
			copy(u.connections[i:], u.connections[i+1:])
			u.connections[len(u.connections)-1] = nil // gc
			u.connections = u.connections[:len(u.connections)-1]
			break
		}
	}
	u.Unlock()
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
