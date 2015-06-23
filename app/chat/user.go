package chat

import (
	"errors"
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

func newUser(id string, name string, location types.LatLng, world types.World) *User {
	u := &User{
		ServerUserJSON: &types.ServerUserJSON{
			BaseServerJSON: &types.BaseServerJSON{
				ID:      id,
				WorldID: world.ID(),
			},
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

func (u *User) SetZone(zone types.Zone) {
	u.Lock()
	defer u.Unlock()

	u.zone = zone
	u.ServerUserJSON.ZoneID = zone.ID()
}

func (u *User) Broadcast(e types.ClientEvent) {
	u.Lock()
	defer u.Unlock()

	for _, connection := range u.connections {
		connection.events <- e
	}
}

func (u *User) Connect() types.Connection {
	connection := newConnection(u)

	u.Lock()
	defer u.Unlock()

	println("connected")
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

func (u *User) ClientJSON() types.ClientJSON {
	return nil
}

func (u *User) ServerJSON() types.ServerJSON {
	return u.ServerUserJSON
}

func (u *User) Update(js types.ServerJSON) error {
	json, ok := js.(*types.ServerUserJSON)
	if !ok {
		return errors.New("Invalid json type.")
	}

	u.Lock()
	defer u.Unlock()
	u.ServerUserJSON = json
	return nil
}
