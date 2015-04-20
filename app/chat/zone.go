package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Zone struct {
	Geohash     string          `json:"geohash"`
	Subscribers []*Subscription `json:"subscribers,omitempty"`
	publish     chan *Event     `json:"-"`
}

func (z *Zone) Type() string {
	return "zone"
}

var zones map[string]*Zone

func FindZone(zone string) (z *Zone, ok bool) {
	println(zone)
	z, ok = zones[zone]
	return
}

func SubscribeToZone(geohash string, user *User) (*Subscription, *Zone) {
	zone, ok := zones[geohash]
	if !ok {
		zone = createZone(geohash)
		zones[geohash] = zone
	}

	subscription := &Subscription{
		User:   user,
		Zone:   zone,
		Events: make(chan *Event, 10),
	}

	zone.Subscribers = append(zone.Subscribers, subscription)
	zone.Publish(NewEvent(&Join{User: user}))
	return subscription, zone
}

func createZone(geohash string) *Zone {
	println(geohash)
	zone := &Zone{
		Geohash:     geohash,
		Subscribers: make([]*Subscription, 0),
		publish:     make(chan *Event, 10),
	}

	go zone.redisSubscribe()
	go zone.run()
	return zone
}

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

func (z *Zone) Publish(event *Event) (*Event, error) {
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
	return z.Publish(NewEvent(m))
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
