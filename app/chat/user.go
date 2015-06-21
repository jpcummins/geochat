package chat

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
	"sync"
	"time"
)

type userJSON struct {
	ID           string       `json:"id"`
	CreatedAt    int          `json:"created_at"`
	LastActivity int          `json:"last_activity"`
	Name         string       `json:"name"`
	Location     types.LatLng `json:"location"`
	ZoneID       string       `json:"zone_id"`
}

type User struct {
	*userJSON
	sync.RWMutex
	zone        types.Zone
	connections []types.Connection
}

func newUser(id string, name string, location types.LatLng) *User {
	u := &User{
		userJSON: &userJSON{
			ID:           id,
			CreatedAt:    int(time.Now().Unix()),
			LastActivity: int(time.Now().Unix()),
			Name:         name,
			Location:     location,
		},
		connections: make([]types.Connection, 0),
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

func (u *User) Location() types.LatLng {
	return u.userJSON.Location
}

func (u *User) Zone() types.Zone {
	return u.zone
}

func (u *User) Broadcast(e types.Event) {
	u.Lock()
	defer u.Unlock()

	for _, connection := range u.connections {
		connection.Events() <- e
	}
}

func (u *User) AddConnection(c types.Connection) {
	u.Lock()
	defer u.Unlock()

	u.connections = append(u.connections, c)
}

func (u *User) RemoveConnection(c types.Connection) {
	u.Lock()
	defer u.Unlock()

	for i, connection := range u.connections {
		if connection == c {
			copy(u.connections[i:], u.connections[i+1:])
			u.connections[len(u.connections)-1] = nil // gc
			u.connections = u.connections[:len(u.connections)-1]
			break
		}
	}
}
