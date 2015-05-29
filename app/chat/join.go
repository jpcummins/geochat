package chat

type Join struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) OnReceive(e *Event) error {
	zone := j.Subscriber.GetZone()
	Subscribers.Set(j.Subscriber)         // Cache Subscriber
	zone.SetSubscription(j.Subscriber)    //
	incrementZoneSubscriptionCounts(zone) // bubble up the count
	zone.broadcastEvent(e)
	return nil
}
