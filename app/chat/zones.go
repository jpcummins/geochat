package chat

import (
    "container/list"
    "time"
    "github.com/garyburd/redigo/redis"
)

var (
    create = make(chan )
)

func zones() {
    zones := list.New()

    for {
        select {
        case zone := <-create:
        case zone := <-destroy:
        case zone := <-rezone:
        case message := <-broadcast:
        }
    }
}

func init() {
    go zones()
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