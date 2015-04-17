package chat

import (
	"container/list"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Zone struct {
	Geohash string

	subscribers *list.List

	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe chan (chan<- Subscription)

	// Send a channel here to unsubscribe.
	unsubscribe chan (<-chan Event)

	// Send events here to publish them.
	publish chan Event
}

func (z *Zone) run() {
	for {
		select {
		case ch := <-z.subscribe:
			subscriber := make(chan Event, 10)
			z.subscribers.PushBack(subscriber)
			ch <- Subscription{subscriber, z}
			subscriber <- *newEvent(z.GetArchive(10))

		case event := <-z.publish:
			for ch := z.subscribers.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}
		case unsub := <-z.unsubscribe:
			for ch := z.subscribers.Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub {
					z.subscribers.Remove(ch)
					break
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
				println("Error:" + err.Error())
				continue
			}
			z.publish <- event
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
	return z.Publish(newEvent(m))
}

func (z *Zone) GetArchive(maxEvents int) *Archive {
	c := pool.Get()
	defer c.Close()

	archiveJson, err := redis.Strings(c.Do("LRANGE", "zone_"+z.Geohash, 0, maxEvents-1))
	if err != nil {
		return nil
	}

	return newArchive(archiveJson)
}
