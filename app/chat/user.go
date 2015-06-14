package chat

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
	"sync"
	"time"
)

type userJSON struct {
	ID           string  `json:"id"`
	CreatedAt    int     `json:"created_at"`
	LastActivity int     `json:"last_activity"`
	Name         string  `json:"name"`
	Lat          float64 `json:"lat"`
	Long         float64 `json:"long"`
}

type User struct {
	*userJSON
	sync.RWMutex
	connections []types.Connection
	world       *World
}

func newUser(lat float64, long float64, name string, id string) *User {
	u := &User{
		userJSON: &userJSON{
			ID:           id,
			CreatedAt:    int(time.Now().Unix()),
			LastActivity: int(time.Now().Unix()),
			Name:         name,
			Lat:          lat,
			Long:         long,
		},
		connections: make([]types.Connection, 0),
		world:       world,
	}
	return u
}

func (u *User) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &u.userJSON)
}

func (u *User) ID() string {
	return u.userJSON.ID
}

func (u *User) Name() string {
	return u.userJSON.Name
}

func (u *User) Broadcast(e types.Event) {
	u.Lock()
	for _, connection := range u.connections {
		connection.Events() <- e
	}
	u.Unlock()
}

func (u *User) AddConnection(c types.Connection) {
	u.Lock()
	u.connections = append(u.connections, c)
	u.Unlock()
}

func (u *User) RemoveConnection(c types.Connection) {
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
