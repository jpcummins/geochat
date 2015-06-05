package chat

type Join struct {
	User *User `json:"user"`
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e *Event) error {
	zone := j.User.GetZone()
	zone.setUser(j.User)
	zone.broadcastEvent(e)

	if _, found := UserCache.Get(j.User.GetID()); !found {
		UserCache.Set(j.User)
	}

	return nil
}
