package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"github.com/TomiHiltunen/geohash-golang"
	"time"
)

type Zone struct {
	*revel.Controller
}

func (c Zone) Lookup(lat float64, long float64) revel.Result {
	return c.RenderText(geohash.EncodeWithPrecision(lat, long, 3))
}

func (c Zone) Message(zone string, text string) revel.Result {

    user, err := chat.GetUser(c.Session["user"])
    if err != nil {
        println("no user")
        return c.RenderError(err)
    }

	z, ok := chat.FindZone(zone)
	if !ok {
		return nil
	}

	event, err := z.SendMessage(user, text)
	if err != nil {
		return c.RenderError(err)
	}

	return c.RenderJson(event)
}

func (c Zone) Zone(zone string) revel.Result {
	box := geohash.Decode(zone)
    user, _ := chat.GetUser(c.Session["user"])

	return c.Render(zone, box, user)
}

func (c Zone) ZoneSocket(zone string, ws *websocket.Conn) revel.Result {

    user, _ := chat.GetUser(c.Session["user"])
    s, z := chat.SubscribeToZone(zone, user)
    defer s.Unsubscribe()

    // Send a heartbeat ping every 30 seconds to keep the Heroku connection alive
	quit := make(chan struct{})
	defer close(quit)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
	    for {
	       select {
	        case <- ticker.C:
	            s.Events <- &chat.Event{Type: "ping", Data: nil, Timestamp: int(time.Now().Unix())}
	        case <- quit:
	            ticker.Stop()
	            return
	        }
	    }
	}()

    // Listen for client disconnects
    go func() {
        var msg string
        for {
            err := websocket.Message.Receive(ws, &msg)
            if err != nil {
            	close(quit)
                s.Unsubscribe()
                return
            }
        }
    }()

    // Send the archive
    if archive, err := z.GetArchive(10); err == nil {
        s.Events <- chat.NewEvent(archive)
    }

	for {
		select {
		case event := <-s.Events:
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
                leave := &chat.Leave{User: user}
                z.Publish(chat.NewEvent(leave))
                close(quit)
				return nil
			}
		}
	}
}
