package chat

import (
	"encoding/json"
	"errors"
)

type joinJSON struct {
	SubscriberID string `json:"subscription"`
}

type Join struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (j *Join) Type() string {
	return "join"
}

func (j *Join) UnmarshalJSON(b []byte) error {
	var js joinJSON
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}

	subscription, found := Subscribers.Get(js.SubscriberID)
	if !found {
		panic(errors.New("Unknown subscription: " + js.SubscriberID))
	}

	j.Subscriber = subscription
	return nil
}

func (j *Join) MarshalJSON() ([]byte, error) {
	joinJSON := &joinJSON{
		SubscriberID: j.Subscriber.GetID(),
	}
	return json.Marshal(joinJSON)
}

func (j *Join) OnReceive(e *Event) error {
	zone := j.Subscriber.GetZone()
	zone.SetSubscription(j.Subscriber)
	incrementZoneSubscriptionCounts(zone) // bubble up the count
	zone.broadcastEvent(e)
	return nil
}
