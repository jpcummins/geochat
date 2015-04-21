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
	return c.RenderText(gh.EncodeWithPrecision(lat, long, 5))
}

func (c Zone) Message(geohash string, text string) revel.Result {
	user, err := chat.GetUser(c.Session["user"])
	if err != nil {
		return c.RenderError(err)
	}

	zone, err := chat.GetOrCreateZone(geohash)
	if err != nil {
		return c.RenderError(err)
	}

	event, err := zone.SendMessage(user, text)
	if err != nil {
		return c.RenderError(err)
	}

	return c.RenderJson(event)
}

func (c Zone) Command(command string, geohash string) revel.Result {
	json, err := chat.ExecuteCommand(command, geohash)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(json)
}

func (c Zone) Zone(geohash string) revel.Result {
	box := gh.Decode(geohash)
	user, _ := chat.GetUser(c.Session["user"])
	return c.Render(geohash, box, user)
}

func (c Zone) ZoneSocket(geohash string, ws *websocket.Conn) revel.Result {
	user, _ := chat.GetUser(c.Session["user"])
	zone, _ := chat.GetOrCreateZone(geohash)
	subscription := zone.Subscribe(user)

	// Listen for client disconnects
	go func() {
		var msg string
		for {
			if websocket.Message.Receive(ws, &msg) != nil {
				zone.Unsubscribe(subscription)
				return
			}
		}
	}()

	// Send zone information
	subscription.Events <- chat.NewEvent(zone)

	// Send the archive
	archive, _ := zone.GetArchive(10)
	subscription.Events <- chat.NewEvent(archive)

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			subscription.Events <- &chat.Event{Type: "ping", Data: nil, Timestamp: int(time.Now().Unix())}
		case event := <-subscription.Events:
			if websocket.JSON.Send(ws, &event) != nil {
				zone.Unsubscribe(subscription)
				return nil
			}
		}
	}
}
