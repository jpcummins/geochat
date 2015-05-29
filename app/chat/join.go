package chat

type Join struct {
	User *User `json:"subscriber"`
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e *Event) error {
	zone := j.User.GetZone()
	UserCache.Set(j.User) // Cache Subscriber
	zone.SetUser(j.User)
	incrementZoneSubscriptionCounts(zone) // bubble up the count
	zone.broadcastEvent(e)
	return nil
}
