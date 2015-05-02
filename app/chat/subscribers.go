package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"strconv"
	"time"
)

type Subscribers struct {
	subscriptions      map[string]*Subscription
	subscribe          chan *Subscription
	unsubscribe        chan *Subscription
	getSubscription    chan (chan interface{})
	publishSubscribe   chan *Subscription
	publishUnsubscribe chan *Subscription
	publishOnline      chan *Subscription
	publishOffline     chan *Subscription
	publishEventToZone chan (chan interface{})
}

func newSubscribers() *Subscribers {
	s := &Subscribers{}

	s.subscriptions = make(map[string]*Subscription)
	s.subscribe = make(chan *Subscription)
	s.unsubscribe = make(chan *Subscription)
	s.getSubscription = make(chan (chan interface{}))

	s.publishSubscribe = make(chan *Subscription)
	s.publishUnsubscribe = make(chan *Subscription)
	s.publishOnline = make(chan *Subscription)
	s.publishOffline = make(chan *Subscription)
	s.publishEventToZone = make(chan (chan interface{}))

	c := connection.Get()
	defer c.Close()
	reply, err := c.Do("HGETALL", "subscribers")
	subscribers, err := redis.StringMap(reply, err)

	if err != nil {
		panic(err)
	}

	go s.redisSubscribe()
	go s.redisPublish()
	go s.manage()

	for _, v := range subscribers {
		var subscription Subscription
		if json.Unmarshal([]byte(v), &subscription) != nil {
			panic("unable to unmarshal subscription")
		}
		s.subscribe <- &subscription
	}

	return s
}

func (s *Subscribers) redisSubscribe() {
	psc := redis.PubSubConn{connection.Get()}
	defer psc.Close()
	psc.Subscribe("subscriptions")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var event Event
			if err := json.Unmarshal(v.Data, &event); err != nil {
				continue
			}

			if event.Type == "join" {
				join := event.Data.(*Join)
				s.subscribe <- join.Subscriber
			}

			if event.Type == "leave" {
				leave := event.Data.(*Leave)
				s.unsubscribe <- leave.Subscriber
			}

			if event.Type == "online" {
				online := event.Data.(*Online)
				online.Subscriber.Activate()
			}

			if event.Type == "offline" {
				offline := event.Data.(*Offline)
				offline.Subscriber.Deactivate()
			}

		}
	}
}

// TODO: Refactor. This can be simplified.
func (s *Subscribers) redisPublish() {
	c := connection.Get()
	defer c.Close()
	for {
		select {
		case subscription := <-s.publishSubscribe:
			join := &Join{subscription}
			eventJson, _ := json.Marshal(NewEvent(join))
			subscriptionJson, _ := json.Marshal(subscription)
			c.Do("PUBLISH", "subscriptions", eventJson)
			c.Do("HSET", "subscribers", subscription.Id, subscriptionJson)
		case subscription := <-s.publishUnsubscribe:
			leave := &Leave{subscription}
			eventJson, _ := json.Marshal(NewEvent(leave))
			subscriptionJson, _ := json.Marshal(subscription)
			c.Do("PUBLISH", "subscriptions", eventJson)
			c.Do("HSET", "subscribers", subscription.Id, subscriptionJson)
		case subscription := <-s.publishOnline:
			online := &Online{subscription}
			eventJson, _ := json.Marshal(NewEvent(online))
			subscriptionJson, _ := json.Marshal(subscription)
			c.Do("PUBLISH", "subscriptions", eventJson)
			c.Do("HSET", "subscribers", subscription.Id, subscriptionJson)
		case subscription := <-s.publishOffline:
			offline := &Offline{subscription}
			eventJson, _ := json.Marshal(NewEvent(offline))
			subscriptionJson, _ := json.Marshal(subscription)
			c.Do("PUBLISH", "subscriptions", eventJson)
			c.Do("HSET", "subscribers", subscription.Id, subscriptionJson)
		}
	}
}

func (s *Subscribers) manage() {
	for {
		select {
		case subscription := <-s.subscribe:
			s.subscriptions[subscription.Id] = subscription
			IncrementZoneSubscriptionCounts(subscription.zone)
		case subscription := <-s.unsubscribe:
			delete(s.subscriptions, subscription.Id)
			DecrementZoneSubscriptionCounts(subscription.zone)
		case ch := <-s.getSubscription:
			id := (<-ch).(string)
			ch <- s.subscriptions[id]
		case ch := <-s.publishEventToZone:
			event := (<-ch).(*Event)
			zone := (<-ch).(*Zone)
			publishEventToZone(event, zone)
			close(ch)
		}
	}
}

func (s *Subscribers) PublishEventToZone(event *Event, zone *Zone) {
	ch := make(chan interface{})
	s.publishEventToZone <- ch
	ch <- event
	ch <- zone
	<-ch
	return
}

func publishEventToZone(event *Event, zone *Zone) {
	for _, subscription := range subscribers.subscriptions {
		if subscription.zone == zone && subscription.Events != nil {
			println("Publishing to user", subscription.User.Name)
			subscription.Events <- event
		}
	}
}

func (s *Subscribers) Get(id string) *Subscription {
	ch := make(chan interface{})
	s.getSubscription <- ch
	ch <- id
	subscription := (<-ch).(*Subscription)
	close(ch)
	return subscription
}

func (s *Subscribers) Add(user *User, zone *Zone) *Subscription {
	subscription := &Subscription{
		Id:       strconv.Itoa(rand.Intn(1000)) + strconv.Itoa(int(time.Now().Unix())),
		User:     user,
		Zonehash: zone.Zonehash,
		Events:   make(chan *Event, 10),
		zone:     zone,
	}
	s.publishSubscribe <- subscription
	zone.Publish(NewEvent(&Join{subscription}))
	return subscription
}

func (s *Subscribers) Remove(subscriber *Subscription) {
	s.publishUnsubscribe <- subscriber
	return
}
