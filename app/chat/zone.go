package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var zones map[string]*Zone

type Zone struct {
	Geohash     string             `json:"geohash"`
	Subscribers []*Subscription    `json:"subscribers"`
	publish     chan *Event        `json:"-"`
	subscribe   chan *Subscription `json:"-"`
	unsubscribe chan *Subscription `json:"-"`
}

func GetOrCreateZone(geohash string) (*Zone, error) {
	zone, found := zones[geohash]

	if !found {
		// TODO: zone validation
		zone := &Zone{
			Geohash:     geohash,
			Subscribers: make([]*Subscription, 0),
			publish:     make(chan *Event, 10),
			subscribe:   make(chan *Subscription, 10),
			unsubscribe: make(chan *Subscription, 10),
		}

		zones[geohash] = zone // unsafe

		go zone.redisSubscribe()
		go zone.run()
		return zone, nil
	}
	return zone, nil
}

func (z *Zone) Type() string {
	return "zone"
}

func (z *Zone) run() {
	for {
		select {
		case subscription := <-z.subscribe:
			z.SendMessage(subscription.User, "a")
			z.Subscribers = append(z.Subscribers, subscription)
			z.SendMessage(subscription.User, "b")
		case subscription := <-z.unsubscribe:
			for i, subscriber := range z.Subscribers {
				if subscriber.Id == subscription.Id {
					copy(z.Subscribers[i:], z.Subscribers[i+1:])
					z.Subscribers[len(z.Subscribers)-1] = nil
					z.Subscribers = z.Subscribers[:len(z.Subscribers)-1]
					break
				}
			}
		case event := <-z.publish:
			for _, subscriber := range z.Subscribers {
				if subscriber != nil {
					subscriber.Events <- event
				}
			}
		}
	}
}

func (z *Zone) Subscribe(user *User) *Subscription {
	subscriber := &Subscription{
		Id:     z.Geohash + user.Id + strconv.Itoa(int(time.Now().Unix())),
		User:   user,
		Zone:   z,
		Events: make(chan *Event, 10),
	}
	z.subscribe <- subscriber
	z.Broadcast(NewEvent(&Join{Subscriber: subscriber}))
	return subscriber
}

func (z *Zone) Unsubscribe(subscriber *Subscription) {
	z.unsubscribe <- subscriber
	z.Broadcast(NewEvent(&Leave{Subscriber: subscriber}))
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
