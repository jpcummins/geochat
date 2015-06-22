package events

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
)

const JoinSeverEvent types.ServerEventType = "join"

type Join struct {
	*types.ServerJoinJSON
	zone types.Zone
	user types.User
}

func NewJoin(zone types.Zone, user types.User) (*Join, error) {
	j := &Join{
		ServerJoinJSON: &types.ServerJoinJSON{
			ZoneID: zone.ID(),
			UserID: user.ID(),
		},
		zone: zone,
		user: user,
	}
	return j, nil
}

func (j *Join) Type() types.ServerEventType {
	return JoinSeverEvent
}

func (j *Join) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &j.ServerJoinJSON); err != nil {
		return err
	}
	return nil
}

func (j *Join) BeforePublish(e types.ServerEvent) error {
	if j.user.Zone() != nil && j.user.Zone() != j.zone {
		// create and publish leave event
	}

	j.zone.AddUser(j.user)
	j.user.SetZone(j.zone)

	if err := e.World().Users().Save(j.user); err != nil {
		return err
	}

	return e.World().Zones().Save(j.zone)
}

func (j *Join) OnReceive(e types.ServerEvent) error {

	if e.World().Zones().FromCache(j.ServerJoinJSON.ZoneID) == nil {
		return nil
	}

	_, err := e.World().Users().FromDB(j.ServerJoinJSON.UserID)
	if err != nil {
		return err
	}

	_, err = e.World().Zones().FromDB(j.ServerJoinJSON.ZoneID)
	if err != nil {
		return nil
	}

	// create and broadcast clientjoinevent
	return nil
}
