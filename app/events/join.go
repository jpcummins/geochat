package events

import (
	"github.com/jpcummins/geochat/app/types"
)

type Join struct {
	world types.World
	zone  types.Zone
	user  types.User
}

func NewJoin(world types.World, zone types.Zone, user types.User) (*Join, error) {
	j := &Join{
		world: world,
		zone:  zone,
		user:  user,
	}
	return j, nil
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) BeforePublish(e types.Event) error {
	if err := j.world.SetUser(j.user); err != nil {
		return err
	}

	j.zone.AddUser(j.user)
	return j.world.SetZone(j.zone)
}

func (j *Join) OnReceive(e types.Event) error {
	return nil
}
