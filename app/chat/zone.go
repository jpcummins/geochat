package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Zone struct {
	Geohash     string          `json:"geohash"`
	Subscribers []*Subscription `json:"subscribers"`
	publish     chan *Event     `json:"-"`
}

func (z *Zone) Type() string {
	return "zone"
}

var zones map[string]*Zone

func (z *Zone) run() {
	for {
		select {
		case event := <-z.publish:
			for _, subscriber := range z.Subscribers {
				if subscriber != nil {
					subscriber.Events <- event
				}
			}
		}
	}
}

func FindZone(geohash string) (z *Zone, ok bool) {
	z, ok = zones[geohash]
	return
}

func GetOrCreateZone(geohash string) (*Zone, error) {
	zone, found := FindZone(geohash)

	if !found {
		return CreateZone(geohash)
	}
	return zone, nil
}

func CreateZone(geohash string) (*Zone, error) {
	// TODO: zone validation
	zone := &Zone{
		Geohash:     geohash,
		Subscribers: make([]*Subscription, 0),
		publish:     make(chan *Event, 10),
	}
	zones[geohash] = zone

	go zone.redisSubscribe()
	go zone.run()
	return zone, nil
}

func (z *Zone) Subscribe(user *User) (*Subscription, error) {
	subscription := CreateSubscription(user, z)
	z.Broadcast(NewEvent(&Join{Subscriber: subscription}))
	z.Subscribers = append(z.Subscribers, subscription)
	return subscription, nil
}

func (z *Zone) Unsubscribe(user *User) {
	for i, subscriber := range z.Subscribers {
		if subscriber.User == user {
			z.Subscribers = append(z.Subscribers[:i], z.Subscribers[i+1:]...)
			z.Broadcast(NewEvent(&Leave{Subscriber: subscriber}))
			break
		}
	}
}

func (z *Zone) Broadcast(event *Event) (*Event, error) {
	c := pool.Get()
	defer c.Close()

	eventJson, err := json.Marshal(event)

	if err != nil {
		println("Unable to marshal event")
		return nil, err
	}

	if _, err := c.Do("PUBLISH", "zone_"+z.Geohash, eventJson); err != nil {
		println("error", err.Error())
		return nil, err
	}

	if _, err := c.Do("LPUSH", "zone_"+z.Geohash, eventJson); err != nil {
		println("error", err.Error())
	}

	return event, nil
}

func (z *Zone) SendMessage(user *User, text string) (*Event, error) {
	m := &Message{User: user, Text: text}
	return z.Broadcast(NewEvent(m))
}

func (z *Zone) GetArchive(maxEvents int) (*Archive, error) {
	c := pool.Get()
	defer c.Close()

	archiveJson, err := redis.Strings(c.Do("LRANGE", "zone_"+z.Geohash, 0, maxEvents-1))
	if err != nil {
		println("unable to get archive:", err.Error())
		return nil, err
	}

	return newArchive(archiveJson), nil
}

func (z *Zone) redisSubscribe() {
	psc := redis.PubSubConn{pool.Get()}
	defer psc.Close()

	psc.Subscribe("zone_" + z.Geohash)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var event Event
			if err := json.Unmarshal(v.Data, &event); err != nil {
				continue
			}
			z.publish <- &event
		}
	}
}
