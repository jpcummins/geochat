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
	world       types.World
	connections []*Connection
}

func newUser(id string, name string, location types.LatLng, world types.World) *User {
	u := &User{
		UserPubSubJSON: &types.UserPubSubJSON{
			ID:   id,
			Name: name,
		},
		location:    location,
		world:       world,
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

func (u *User) ZoneID() string {
	return u.UserPubSubJSON.ZoneID
}

func (u *User) SetZoneID(id string) {
	u.Lock()
	defer u.Unlock()
	u.UserPubSubJSON.ZoneID = id
}

func (u *User) Broadcast(data types.BroadcastEventData) {
	event := broadcast.NewEvent(generateEventID(), data)

	u.Lock()
	defer u.Unlock()
	for _, connection := range u.connections {
		connection.events <- event
	}
}

func (u *User) ExecuteCommand(command string, args string) error {
	return commands.Execute(command, args, u, u.world)
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
			copy(u.connections[i:], u.connections[i+1:])
			u.connections[len(u.connections)-1] = nil // gc
			u.connections = u.connections[:len(u.connections)-1]
			close(connection.events)
			break
		}
	}
}

func (u *User) BroadcastJSON() interface{} {
	return &types.UserBroadcastJSON{
		ID:       u.ID(),
		Name:     u.Name(),
		Location: u.Location().BroadcastJSON().(*types.LatLngJSON),
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
