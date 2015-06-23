package events

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
)

const JoinSeverEvent types.ServerEventType = "join"

type serverJoin struct {
	*types.ServerJoinJSON
	zone types.Zone
	user types.User
}

func ServerJoin(zone types.Zone, user types.User) *serverJoin {
	j := &serverJoin{
		ServerJoinJSON: &types.ServerJoinJSON{
			Zone: zone.ServerJSON(),
			User: user.ServerJSON(),
		},
		zone: zone,
		user: user,
	}
	return j
}

func (j *serverJoin) Type() types.ServerEventType {
	return JoinSeverEvent
}

func (j *serverJoin) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &j.ServerJoinJSON); err != nil {
		return err
	}
	return nil
}

func (j *serverJoin) BeforePublish(e types.ServerEvent) error {
	if j.user.Zone() != nil && j.user.Zone() != j.zone {
		// create and publish leave event
	}
	println("before publish")

	j.zone.AddUser(j.user)
	j.user.SetZone(j.zone)

	if err := e.World().Users().Save(j.user); err != nil {
		return err
	}

	return e.World().Zones().Save(j.zone)
}

func (j *serverJoin) OnReceive(e types.ServerEvent) error {

	println("1")
	zone := e.World().Zones().FromCache(j.ServerJoinJSON.Zone.Key())
	if zone == nil {
		return nil
	}

	println("2")
	user := e.World().Users().FromCache(j.ServerJoinJSON.User.Key())
	if user == nil {
		var err error
		user, err = e.World().Users().FromDB(j.ServerJoinJSON.User.Key())
		if err != nil {
			return err
		}
	} else {
		println("3")
		user.Update(j.ServerJoinJSON.User)
	}

	println("4")
	zone.Update(j.ServerJoinJSON.User)

	println("5")
	join := e.World().NewClientEvent(ClientJoin(e.ID(), user, zone))
	println("6")
	zone.Broadcast(join)
	println("7")
	return nil
}

const joinClientEvent types.ClientEventType = "join"

type clientJoin struct {
	*types.BaseClientJSON
	User types.ClientJSON `json:"user"`
	Zone types.ClientJSON `json:"zone"`
}

func ClientJoin(id string, user types.User, zone types.Zone) *clientJoin {
	return &clientJoin{
		BaseClientJSON: &types.BaseClientJSON{
			ID: id,
		},
		User: user.ClientJSON(),
		Zone: zone.ClientJSON(),
	}
}

func (c *clientJoin) Type() types.ClientEventType {
	return joinClientEvent
}

func (c *clientJoin) BeforeBroadcast(data types.ClientEvent) error {
	return nil
}
