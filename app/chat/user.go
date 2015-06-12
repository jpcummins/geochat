package chat

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
	"sync"
	"time"
)

type userJSON struct {
	Id           string  `json:"id"`
	CreatedAt    int     `json:"created_at"`
	LastActivity int     `json:"last_activity"`
	Name         string  `json:"name"`
	Lat          float64 `json:"lat"`
	Long         float64 `json:"long"`
}

type User struct {
	*userJSON
	sync.RWMutex
	connections []*Connection
	world       *World
}

func newUser(world *World, lat float64, long float64, name string, id string) (*User, error) {
	u := &User{
		userJSON: &userJSON{
			Id:           id,
			CreatedAt:    int(time.Now().Unix()),
			LastActivity: int(time.Now().Unix()),
			Name:         name,
			Lat:          lat,
			Long:         long,
		},
		connections: make([]*Connection, 0),
		world:       world,
	}
	err := world.cache.SetUser(u)
	return u, err
}

func (u *User) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &u.userJSON)
}

func (u *User) ID() string {
	return u.userJSON.Id
}

func (u *User) Name() string {
	return u.userJSON.Name
}

func (u *User) NewConnection() (types.Connection, error) {
	c := newConnection(u)
	u.Lock()
	u.connections = append(u.connections, c)
	u.Unlock()
	return c, nil
}

func (u *User) Disconnect(c types.Connection) error {
	close(c.Events())
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
	return nil
}

// // ExecuteCommand allows certain subscribers to issue administrative commands.
// func (u *User) ExecuteCommand(command string) (string, error) {
// 	args := strings.Split(command, " ")
// 	if len(args) == 0 || commands[args[0]] == nil {
// 		output, err := json.Marshal(commands)
// 		return string(output[:]), err
// 	}
// 	return commands[args[0]].execute(args[1:], u)
// }
//
// // UpdateLastActiveTime sets the last active time for the subscriber
// func (u *User) UpdateLastActiveTime() {
// 	u.lastActivity = int(time.Now().Unix())
// 	u.zone.world.users.Set(u)
// }
