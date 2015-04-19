package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"github.com/TomiHiltunen/geohash-golang"
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

    // Listen for client disconnects
    go func() {
        var msg string
        for {
            err := websocket.Message.Receive(ws, &msg)
            if err != nil {
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
				return nil
			}
		}
	}
}
