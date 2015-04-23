package chat

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Zone struct {
	Geohash     string                    `json:"geohash"`
	Subhash     int                       `json:"subhash"`
	Parent      *Zone                     `json:"-"`
	Left        *Zone                     `json:"-"`
	Right       *Zone                     `json:"-"`
	Count       int32                     `json:"-"`
	Subscribers []*Subscription           `json:"subscribers"`
	publish     chan *Event               `json:"-"`
	broadcast   chan *Event               `json:"-"`
	subscribe   chan (chan *Subscription) `json:"-"`
	unsubscribe chan (chan *Subscription) `json:"-"`
}

var world = createZone("", 0, nil)

func createZone(geohash string, subhash int, parent *Zone) *Zone {
	zone := &Zone{
		Geohash:     geohash,
		Subhash:     subhash,
		Parent:      parent,
		Subscribers: make([]*Subscription, 0),
		publish:     make(chan *Event), // unbuffered communication with web sockets
		broadcast:   make(chan *Event), // unbuffered communication with redis
		subscribe:   make(chan (chan *Subscription), 10),
		unsubscribe: make(chan (chan *Subscription), 10),
	}
	go zone.redisSubscribe()      // subscribes to redis channel and publishes events
	go zone.redisPublish()        // publishes broadcast events to redis channel
	go zone.manageSubscriptions() // handles communication with zone subscribers
	return zone
}

func (z *Zone) Type() string {
	return "zone"
}

func FindAvailableZone(geohash string) (*Zone, error) {
	return findChatZone(world, geohash)
}

func findChatZone(root *Zone, geohash string) (*Zone, error) {
	if root.Left == nil && root.Right == nil {
		if root.Subhash < 15 {
			root.Left = createZone(root.Geohash, (root.Subhash*2)+1, root)
			root.Right = createZone(root.Geohash, (root.Subhash*2)+2, root)
		} else {
			geohashmap := "0123456789bcdefghjkmnprstuvwxyz"
			root.Left = createZone(root.Geohash+string(geohashmap[(root.Subhash*2)-30]), 0, root)
			root.Right = createZone(root.Geohash+string(geohashmap[(root.Subhash*2)-29]), 0, root)
		}
	}

	if root.Count < maxRoomSize {
		println(root.Geohash)
		return root, nil
	}

	suffix := strings.TrimPrefix(geohash, root.Geohash)

	// edge case - Zone for the specified geohash is full.
	if len(suffix) == 0 {
		return root, errors.New("Room full")
	}

	if suffix[0] < 'g' {
		return findChatZone(root.Left, geohash)
	} else {
		return findChatZone(root.Right, geohash)
	}
}

func (z *Zone) GetZonehash() string {
	println(z.Geohash)
	return z.Geohash + ":" + strconv.Itoa(z.Subhash)
}

func (z *Zone) manageSubscriptions() {
	for {
		select {
		case ch := <-z.subscribe:
			subscription := <-ch
			z.Subscribers = append(z.Subscribers, subscription)

			// update counts
			for zone := z; zone != nil; zone = zone.Parent {
				atomic.AddInt32(&zone.Count, 1)
				runtime.Gosched()
			}

			subscriptions[subscription.Id] = subscription // unsafe

			ch <- subscription
		case ch := <-z.unsubscribe:
			subscription := <-ch
			for i, subscriber := range z.Subscribers {
				if subscriber.Id == subscription.Id {

					// decrement counts
					for zone := z; zone != nil; zone = zone.Parent {
						atomic.AddInt32(&zone.Count, -1)
						runtime.Gosched()
					}

					// TODO: delete subscription after some TTL

					copy(z.Subscribers[i:], z.Subscribers[i+1:])
					z.Subscribers[len(z.Subscribers)-1] = nil
					z.Subscribers = z.Subscribers[:len(z.Subscribers)-1]
					ch <- subscription
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
	var newSubscription = &Subscription{
		Id:     strconv.Itoa(rand.Intn(1000)) + strconv.Itoa(int(time.Now().Unix())),
		User:   user,
		Zone:   z,
		Events: make(chan *Event, 10),
	}
	req := make(chan *Subscription)
	z.subscribe <- req     // add channel to queue
	req <- newSubscription // when ready, pass the subscription
	subscription := <-req  // wait for processing to finish
	z.Broadcast(NewEvent(&Join{Subscriber: subscription}))
	return subscription
}

func (z *Zone) Unsubscribe(subscriber *Subscription) {
	req := make(chan *Subscription)
	z.unsubscribe <- req
	req <- subscriber
	<-req
	z.Broadcast(NewEvent(&Leave{Subscriber: subscriber}))
}

func (z *Zone) Broadcast(event *Event) {
	z.broadcast <- event
}

func (z *Zone) SendMessage(user *User, text string) *Event {
	m := &Message{User: user, Text: text}
	e := NewEvent(m)
	z.Broadcast(e)
	return e
}

func (z *Zone) GetArchive(maxEvents int) (*Archive, error) {
	c := pool.Get()
	defer c.Close()

	archiveJson, err := redis.Strings(c.Do("LRANGE", "zone_"+z.GetZonehash(), 0, maxEvents-1))
	if err != nil {
		println("unable to get archive:", err.Error())
		return nil, err
	}

	return newArchive(archiveJson), nil
}

func (z *Zone) redisSubscribe() {
	psc := redis.PubSubConn{pool.Get()}
	defer psc.Close()

	psc.Subscribe("zone_" + z.GetZonehash())
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

func (z *Zone) redisPublish() {
	c := pool.Get()
	defer c.Close()

	for {
		select {
		case event := <-z.broadcast:
			eventJson, _ := json.Marshal(event)
			c.Do("LPUSH", "zone_"+z.GetZonehash(), eventJson)
			c.Do("PUBLISH", "zone_"+z.GetZonehash(), eventJson)
		}
	}
}
