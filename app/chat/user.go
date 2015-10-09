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

func newUser(id string, world types.World) *User {
	u := &User{
		UserPubSubJSON: &types.UserPubSubJSON{
			ID: id,
		},
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

func (u *User) SetName(name string) {
	u.UserPubSubJSON.Name = name
}

func (u *User) FirstName() string {
	return u.UserPubSubJSON.FirstName
}

func (u *User) SetFirstName(firstName string) {
	u.UserPubSubJSON.FirstName = firstName
}

func (u *User) LastName() string {
	return u.UserPubSubJSON.LastName
}

func (u *User) SetLastName(lastName string) {
	u.UserPubSubJSON.LastName = lastName
}

func (u *User) Timezone() float64 {
	return u.UserPubSubJSON.Timezone
}

func (u *User) SetTimezone(timezone float64) {
	u.UserPubSubJSON.Timezone = timezone
}

func (u *User) Email() string {
	return u.UserPubSubJSON.Email
}

func (u *User) SetEmail(email string) {
	u.UserPubSubJSON.Email = email
}

func (u *User) Location() types.LatLng {
	return u.location
}

func (u *User) SetLocation(lat float64, lng float64) {
	u.location = newLatLng(lat, lng)
}

func (u *User) Locality() string {
	return u.UserPubSubJSON.Locality
}

func (u *User) SetLocality(locality string) {
	u.UserPubSubJSON.Locality = locality
}

func (u *User) ZoneID() string {
	return u.UserPubSubJSON.ZoneID
}

func (u *User) SetZoneID(id string) {
	u.UserPubSubJSON.ZoneID = id
}

func (u *User) FBID() string {
	return u.UserPubSubJSON.FBID
}

func (u *User) SetFBID(id string) {
	u.UserPubSubJSON.FBID = id
}

func (u *User) FBAccessToken() string {
	return u.UserPubSubJSON.FBAccessToken
}

func (u *User) SetFBAccessToken(accessToken string) {
	u.UserPubSubJSON.FBAccessToken = accessToken
}

func (u *User) FBPictureURL() string {
	return u.UserPubSubJSON.FBPictureURL
}

func (u *User) SetFBPictureURL(url string) {
	u.UserPubSubJSON.FBPictureURL = url
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
	user := &types.UserBroadcastJSON{
		ID:           u.ID(),
		FirstName:    u.FirstName(),
		FBPictureURL: u.FBPictureURL(),
		Locality:     u.Locality(),
	}

	if latlng := u.Location(); latlng != nil {
		user.Location = latlng.BroadcastJSON().(*types.LatLngJSON)
	}

	return user
}

func (u *User) PubSubJSON() types.PubSubJSON {

	if u.location != nil {
		u.UserPubSubJSON.Location = u.location.PubSubJSON().(*types.LatLngJSON)
	}

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

	if json.Location != nil {
		u.location = newLatLng(json.Location.Lat, json.Location.Lng)
	}

	return nil
}
