package chat

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
	"sync"
	"time"
)

type User struct {
	*types.ServerUserJSON
	sync.RWMutex
	zone        types.Zone
	connections []*Connection
}

func newUser(id string, name string, location types.LatLng) *User {
	u := &User{
		ServerUserJSON: &types.ServerUserJSON{
			ID:           id,
			CreatedAt:    int(time.Now().Unix()),
			LastActivity: int(time.Now().Unix()),
			Name:         name,
			Location:     location,
		},
		connections: make([]*Connection, 0),
	}
	return u
}

func (u *User) ID() string {
	return u.ServerUserJSON.ID
}

func (u *User) Name() string {
	return u.ServerUserJSON.Name
}

func (u *User) Location() types.LatLng {
	return u.ServerUserJSON.Location
}

func (u *User) Zone() types.Zone {
	return u.zone
}

func (u *User) Broadcast(e types.ClientEvent) {
	u.Lock()
	defer u.Unlock()

	for _, connection := range u.connections {
		connection.Events() <- e
	}
}

func (u *User) Connect() types.Connection {
	connection := newConnection(u)

	u.Lock()
	defer u.Unlock()

	u.connections = append(u.connections, connection)
	return connection
}

func (u *User) Disconnect(c types.Connection) {
	u.Lock()
	defer u.Unlock()

	for i, connection := range u.connections {
		if connection == c {
			close(connection.events)
			copy(u.connections[i:], u.connections[i+1:])
			u.connections[len(u.connections)-1] = nil // gc
			u.connections = u.connections[:len(u.connections)-1]
			break
		}
	}
}

func (u *User) ClientJSON() ([]byte, error) {
	return nil, nil
}

func (u *User) ServerJSON() ([]byte, error) {
	return json.Marshal(u.ServerUserJSON)
}

func (u *User) Update(js *types.ServerUserJSON) error {
	u.Lock()
	defer u.Unlock()

	u.ServerUserJSON = js
	return nil
}
