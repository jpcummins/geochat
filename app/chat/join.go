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
	UserCache.cacheSet(j.User)
	return nil
}
