package chat

import (
	"encoding/json"
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var geohashmap = "0123456789bcdefghjkmnpqrstuvwxyz"

type Zone struct {
	Zonehash    string                    `json:"zonehash"`
	geohash     string                    `json:"-"`
	from        byte                      `json:"-"`
	to          byte                      `json:"-"`
	parent      *Zone                     `json:"-"`
	left        *Zone                     `json:"-"`
	right       *Zone                     `json:"-"`
	Count       int32                     `json:"-"`
	Subscribers []*Subscription           `json:"subscribers"`
	publish     chan *Event               `json:"-"`
	broadcast   chan *Event               `json:"-"`
	subscribe   chan (chan *Subscription) `json:"-"`
	unsubscribe chan (chan *Subscription) `json:"-"`
}

func createZone(geohash string, from byte, to byte, parent *Zone) *Zone {
	zone := &Zone{
		Zonehash: geohash + ":" + string(from) + string(to),
		geohash:  geohash,
		from:     from,
		to:       to,
		parent:   parent,
	}
	return zone
}

func (z *Zone) init() {
	z.Subscribers = make([]*Subscription, 0)
	z.publish = make(chan *Event)   // unbuffered communication with web sockets
	z.broadcast = make(chan *Event) // unbuffered communication with redis
	z.subscribe = make(chan (chan *Subscription), 10)
	z.unsubscribe = make(chan (chan *Subscription), 10)
	go z.redisSubscribe()      // subscribes to redis channel and publishes events
	go z.redisPublish()        // publishes broadcast events to redis channel
	go z.manageSubscriptions() // handles communication with zone subscribers
}

func (z *Zone) createChildZones() {
	from_i := strings.Index(geohashmap, string(z.from))
	to_i := strings.Index(geohashmap, string(z.to))

	if to_i-from_i > 1 {
		split := (to_i - from_i) / 2
		z.left = createZone(z.geohash, z.from, geohashmap[from_i+split], z)
		z.right = createZone(z.geohash, geohashmap[from_i+split+1], z.to, z)
	} else {
		z.left = createZone(z.geohash+string(z.from), '0', 'z', z)
		z.right = createZone(z.geohash+string(z.to), '0', 'z', z)
	}
}

func (z *Zone) Type() string {
	return "zone"
}

func FindAvailableZone(lat float64, long float64) (*Zone, error) {
	geohash := gh.EncodeWithPrecision(lat, long, 6)
	return findChatZone(world, geohash)
}

func findChatZone(root *Zone, geohash string) (*Zone, error) {
	if root.left == nil && root.right == nil {
		root.createChildZones()
	}

	if root.Count < maxRoomSize {
		return root, nil
	}

	suffix := strings.TrimPrefix(geohash, root.geohash)

	if len(suffix) == 0 {
		return root, errors.New("Room full")
	}

	// This is gross. Like, really gross.
	if root.geohash == root.right.geohash {
		if suffix[0] < root.right.from {
			return findChatZone(root.left, geohash)
		} else {
			return findChatZone(root.right, geohash)
		}
	} else {
		if suffix[0] < root.right.geohash[len(root.right.geohash)-1] {
			return findChatZone(root.left, geohash)
		} else {
			return findChatZone(root.right, geohash)
		}
	}

}

func (z *Zone) manageSubscriptions() {
	for {
		select {
		case ch := <-z.subscribe:
			subscription := <-ch
			z.Subscribers = append(z.Subscribers, subscription)

			// update counts
			for zone := z; zone != nil; zone = zone.parent {
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
					for zone := z; zone != nil; zone = zone.parent {
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

	archiveJson, err := redis.Strings(c.Do("LRANGE", "zone_"+z.Zonehash, 0, maxEvents-1))
	if err != nil {
		println("unable to get archive:", err.Error())
		return nil, err
	}

	return newArchive(archiveJson), nil
}

func (z *Zone) GetBoundries() *gh.BoundingBox {
	return gh.Decode(z.geohash)
}

func (z *Zone) redisSubscribe() {
	psc := redis.PubSubConn{pool.Get()}
	defer psc.Close()

	psc.Subscribe("zone_" + z.Zonehash)
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
			c.Do("LPUSH", "zone_"+z.Zonehash, eventJson)
			c.Do("PUBLISH", "zone_"+z.Zonehash, eventJson)
		}
	}
}
