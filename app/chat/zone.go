package chat

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	gh "github.com/TomiHiltunen/geohash-golang"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"math"
)

type Zone struct {
	Zonehash    string                    `json:"zonehash"`
	geohash     string                    `json:"-"`
	subhash     int                       `json:"-"`
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

var world = createZone("", 1, nil)
var geohashmap = "0123456789bcdefghjkmnpqrstuvwxyz"

func createZone(geohash string, subhash int, parent *Zone) *Zone {
	zone := &Zone{
		Zonehash:    geohash + ":" + strconv.Itoa(subhash),
		geohash:     geohash,
		subhash:     subhash,
		parent:      parent,
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

func (z *Zone) createChildZones() {
	if z.subhash < 16 {
		z.left = createZone(z.geohash, (z.subhash*2), z)
		z.right = createZone(z.geohash, (z.subhash*2)+1, z)
	} else {
		z.left = createZone(z.geohash+string(geohashmap[(z.subhash*2)-32]), 1, z)
		z.right = createZone(z.geohash+string(geohashmap[(z.subhash*2)-31]), 1, z)
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
	if (root.left == nil && root.right == nil) {
		root.createChildZones()
	}

	if root.Count < maxRoomSize {
		return root, nil
	}

	suffix := strings.TrimPrefix(geohash, root.geohash)

	// edge case - Zone for the specified geohash is full.
	if len(suffix) == 0 {
		return root, errors.New("Room full")
	}

	// GROSS. I'm not a mathematician nor am I an algorithms expert.
	// I'm sorry if this makes your eyes bleed.
	l := math.Pow(2, math.Ceil(math.Log2(float64(root.right.subhash) + 1)) - 1)
	d := 32 / l
	i := int(d * (float64(root.right.subhash) - l))

	if suffix[0] < geohashmap[i] {
		return findChatZone(root.left, geohash)
	} else {
		return findChatZone(root.right, geohash)
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
