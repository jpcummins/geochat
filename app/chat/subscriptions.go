package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"sync"
)

type Subscriptions struct {
	sync.RWMutex
	subscriptions map[string]*Subscription
}

func NewSubscriptions() *Subscriptions {
	s := &Subscriptions{
		subscriptions: make(map[string]*Subscription),
	}
	return s
}

func (s *Subscriptions) Get(id string) (*Subscription, bool) {
	if subscription, ok := s.cacheGet(id); ok {
		return subscription, ok
	}

	if subscription, ok := s.redisGet(id); ok {
		s.cacheSet(subscription)
		return subscription, ok
	}

	return nil, false
}

func (s *Subscriptions) Set(subscription *Subscription) {
	s.cacheSet(subscription)
	s.redisSet(subscription)
}

func (s *Subscriptions) cacheGet(id string) (*Subscription, bool) {
	s.RLock()
	subscription, found := s.subscriptions[id]
	s.RUnlock()
	return subscription, found
}

func (s *Subscriptions) cacheSet(subscription *Subscription) {
	s.Lock()
	s.subscriptions[subscription.GetID()] = subscription
	s.Unlock()
	subscription.GetZone().SetSubscription(subscription)
}

func (s *Subscriptions) redisGet(id string) (*Subscription, bool) {
	c := connection.Get()
	defer c.Close()

	subscriptionJSON, err := redis.String(c.Do("HGET", "subscriptions", id))
	if err != nil {
		return nil, false
	}

	var subscription Subscription
	if err := json.Unmarshal([]byte(subscriptionJSON), &subscription); err != nil {
		return nil, false
	}
	return &subscription, true
}

func (s *Subscriptions) redisSet(subscription *Subscription) {
	c := connection.Get()
	defer c.Close()
	eventJSON, err := json.Marshal(subscription)
	if err != nil {
		return
	}
	c.Do("HSET", "subscriptions", subscription.id, eventJSON)
}
