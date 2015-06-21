package events

import (
	"encoding/json"
	"github.com/jpcummins/geochat/app/types"
)

type userJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type joinJSON struct {
	WorldID  string    `json:"world_id"`
	ZoneID   string    `json:"zone_id"`
	UserJSON *userJSON `json:"user"`
}

type Join struct {
	*joinJSON
	zone types.Zone
	user types.User
}

func NewJoin(zone types.Zone, user types.User) (*Join, error) {
	j := &Join{
		zone: zone,
		user: user,
		joinJSON: &joinJSON{
			ZoneID: zone.ID(),
			UserJSON: &userJSON{
				ID:   user.ID(),
				Name: user.Name(),
			},
		},
	}
	return j, nil
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &j.joinJSON); err != nil {
		return err
	}
	return nil
}

func (j *Join) BeforePublish(e types.Event) error {
	if err := e.World().Users().SetUser(j.user); err != nil {
		return err
	}

	j.zone.AddUser(j.user)
	return e.World().Zones().SetZone(j.zone)
}

func (j *Join) OnReceive(e types.Event) error {
	_, err := e.World().Users().UpdateUser(j.joinJSON.UserJSON.ID)
	if err != nil {
		return err
	}

	zone, err := e.World().Zones().UpdateZone(j.joinJSON.ZoneID)
	if err != nil {
		return nil
	}

	zone.Broadcast(e)
	return nil
}
