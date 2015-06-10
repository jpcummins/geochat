package chat

import (
	"encoding/json"
)

type Join struct {
	user *User
}

func (j *Join) UnmarshalJSON(b []byte) error {
	var js userJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	if user, found := UserCache.Get(js.ID); found {
		j.user = user
		return nil
	}

	j.user = &User{}
	return j.user.UnmarshalJSON(b)
}

func (j *Join) MarshalJSON() ([]byte, error) {
	return j.user.MarshalJSON()
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e *Event, z *Zone) error {
	z.addUser(j.user)
	z.broadcastEvent(e)
	return nil
}
