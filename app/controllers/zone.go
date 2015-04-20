package controllers

import (
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"time"
)

type Zone struct {
	*revel.Controller
}

func (c Zone) Lookup(lat float64, long float64) revel.Result {
	return c.RenderText(gh.EncodeWithPrecision(lat, long, 3))
}

func (c Zone) Message(geohash string, text string) revel.Result {
	user, err := chat.GetUser(c.Session["user"])
	if err != nil {
		return c.RenderError(err)
	}

	z, ok := chat.FindZone(geohash)
	if !ok {
		return nil
	}

	event, err := z.SendMessage(user, text)
	if err != nil {
		return c.RenderError(err)
	}

	return c.RenderJson(event)
}

func (c Zone) Zone(geohash string) revel.Result {
	box := gh.Decode(geohash)
	user, _ := chat.GetUser(c.Session["user"])
	return c.Render(geohash, box, user)
}

func (c Zone) ZoneSocket(geohash string, ws *websocket.Conn) revel.Result {
	user, _ := chat.GetUser(c.Session["user"])
	s, zone := chat.SubscribeToZone(geohash, user)
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

	// Send zone information
	s.Events <- chat.NewEvent(zone)

	// Send the archive
	if archive, err := zone.GetArchive(10); err == nil {
		s.Events <- chat.NewEvent(archive)
	}

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			s.Events <- &chat.Event{Type: "ping", Data: nil, Timestamp: int(time.Now().Unix())}
			continue
		case event := <-s.Events:
			if err := websocket.JSON.Send(ws, &event); err != nil {
				// They disconnected.
				leave := &chat.Leave{User: user}
				zone.Publish(chat.NewEvent(leave))
				return nil
			}
		}
	}
}
