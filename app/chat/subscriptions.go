package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Subscriptions struct{}

func (s *Subscriptions) Get(id string) (subscription *Subscription, found bool) {
	c := connection.Get()
	defer c.Close()

	subscriptionJSON, err := redis.String(c.Do("HGET", "subscriptions", id))
	if err != nil {
		return nil, false
	}

	if err := json.Unmarshal([]byte(subscriptionJSON), &subscription); err != nil {
		return nil, false
	}

	return subscription, true
}

func (s *Subscriptions) Set(subscription *Subscription) {
	c := connection.Get()
	defer c.Close()
	eventJSON, err := json.Marshal(subscription)
	if err != nil {
		return
	}
	c.Do("HSET", "subscriptions", subscription.id, eventJSON)
}
