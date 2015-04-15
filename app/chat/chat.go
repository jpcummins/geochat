package chat

import (
    "time"
    "github.com/garyburd/redigo/redis"
    "container/list"
)

type Event struct {
    Type      string // "join", "leave", or "message"
    User      string
    Timestamp int    // Unix timestamp 
    Text      string // What the user said (if Type == "message")
}

type Subscription struct {
    New <-chan Event
    Zone *Zone
}

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

var zones map[string]*Zone

func (z *Zone) run() {
    for {
        select {
        case ch := <-z.subscribe:
            subscriber := make(chan Event, 10)
            z.subscribers.PushBack(subscriber)
            ch <- Subscription{subscriber, z}

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

func (z *Zone) Say(user string, text string) *Event {
    e := newEvent("message", user, text)
    z.publish <- e
    return &e
}

func newEvent(typ, user, msg string) Event {
    return Event{typ, user, int(time.Now().Unix()), msg}
}

func FindZone(zone string) (z *Zone, ok bool)  {
    z, ok = zones[zone]
    return
}

func createZone(geohash string) *Zone {
    zone := &Zone{
        Geohash: geohash, 
        subscribers: list.New(),
        subscribe: make(chan (chan<- Subscription), 10),
        unsubscribe: make(chan (<-chan Event), 10),
        publish: make(chan Event, 10),
    }

    go zone.run()
    return zone
}

func (s Subscription) Unsubscribe() {
    s.Zone.unsubscribe <- s.New
}

func SubscribeToZone(geohash string) Subscription {
    zone, ok := zones[geohash] 
    if !ok {
        zone = createZone(geohash)
        zones[geohash] = zone
    }

    subscription := make(chan Subscription)
    zone.subscribe <- subscription
    return <-subscription
}

func createPool(server string) *redis.Pool {
    return &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}

var (
    pool *redis.Pool
)

func init() {
    pool = createPool(":6379")
    zones = make(map[string]*Zone)
}