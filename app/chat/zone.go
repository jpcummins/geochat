package chat

import (
	"container/list"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Zone struct {
	Geohash string

	subscribers *list.List

	// Send events here to publish them.
	publish chan *Event
}

func (z *Zone) run() {
	for {
		select {
		case event := <-z.publish:
			for s := z.subscribers.Front(); s != nil; s = s.Next() {
				subscriber := s.Value.(*Subscription)
				subscriber.Events <- event
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
				println("Error:" + err.Error())
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