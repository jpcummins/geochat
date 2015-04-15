package main

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
    "time"
)

type Zone struct {
    Geohash string
}

func Create(geohash string) *Zone {
    zone := &Zone{Geohash: geohash}

    // Tell the world
    c := pool.Get()
    defer pool.Close()

    if _, err := c.Do("PUBLISH", "zone", zone.Geohash); err != nil {
        panic(err)
    }

    println("published")

    return zone
}

func zones() {
    c := pool.Get()
    defer c.Close()

    psc := redis.PubSubConn{c}
    if err := psc.Subscribe("zone"); err != nil {
        panic(err)
    }

    for {
        switch v := psc.Receive().(type) {
        case redis.Message:
            fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
        case redis.Subscription:
            fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
        case error:
            println("error")
        }
    }
}

func newPool(server string) *redis.Pool {
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

func main() {
    pool = newPool(":6379")

    go zones()

    timer := time.NewTimer(time.Second * 2)
    <- timer.C
    Create("abc")
}


/*

Notes

global:
rezone
destroy

per zone:
subscribe
unsubscribe
publish


*/