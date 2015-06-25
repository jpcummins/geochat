package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/broadcast"
	"github.com/jpcummins/geochat/app/commands"
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type User struct {
	sync.RWMutex
	types.PubSubSerializable
	types.BroadcastSerializable
	*types.UserPubSubJSON
	location    types.LatLng
	zone        types.Zone
	connections []*Connection
}

func newUser(id string, name string, location types.LatLng, world types.World) *User {
	u := &User{
		UserPubSubJSON: &types.UserPubSubJSON{
			ID:   id,
			Name: name,
		},
		location:    location,
		connections: make([]*Connection, 0),
	}
	return u
}

func (u *User) ID() string {
	return u.UserPubSubJSON.ID
}

func (u *User) Name() string {
	return u.UserPubSubJSON.Name
}

func (u *User) Location() types.LatLng {
	return u.location
}

func (u *User) Zone() types.Zone {
	return u.zone
}

func (u *User) SetZone(zone types.Zone) {
	u.Lock()
	defer u.Unlock()

	u.zone = zone
	if zone != nil {
		u.UserPubSubJSON.ZoneID = zone.ID()
	}
}

func (u *User) Broadcast(data types.BroadcastEventData) {
	event := broadcast.NewEvent(generateEventID(), data)

	if ok, err := data.BeforeBroadcastToUser(u, event); !ok || err != nil {
		return
	}

	u.Lock()
	defer u.Unlock()
	for _, connection := range u.connections {
		connection.events <- event
	}
}

func (u *User) ExecuteCommand(command string, args string) error {
	return commands.Execute(command, args, u)
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

func (u *User) BroadcastJSON() interface{} {
	return &types.UserBroadcastJSON{
		ID:   u.ID(),
		Name: u.Name(),
	}
}

func (u *User) PubSubJSON() types.PubSubJSON {
	u.UserPubSubJSON.Location = u.location.PubSubJSON().(*types.LatLngJSON)
	return u.UserPubSubJSON
}

func (u *User) Update(js types.PubSubJSON) error {
	json, ok := js.(*types.UserPubSubJSON)
	if !ok {
		return errors.New("Unable to serialize to UserPubSubJSON.")
	}

	u.Lock()
	defer u.Unlock()
	u.UserPubSubJSON = json
	u.location = newLatLng(json.Location.Lat, json.Location.Lng)
	return nil
}
