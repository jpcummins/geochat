package controllers

import (
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"time"
	"strconv"
	"errors"
)

type Zone struct {
	*revel.Controller
}

func (c Zone) Lookup(lat float64, long float64) revel.Result {
	geohash := gh.EncodeWithPrecision(lat, long, 8)

	type LookupResponse struct {
		Geohash string `json:"geohash"`
		Zonehash string `json:"zonehash"`
	}

	zone, _ := chat.FindAvailableZone(geohash)
	resp := &LookupResponse{geohash, zone.Geohash + ":" + strconv.Itoa(zone.Subhash)}
	return c.RenderJson(resp)
}

func (c Zone) Message(geohash string, text string) revel.Result {
	user, err := chat.GetUser(c.Session["user"])
	if err != nil {
		return c.RenderError(err)
	}

	// TODO: gross. This is not right.
	zone, err := chat.FindAvailableZone(geohash)
	if (zone.Geohash != geohash) {
		return c.RenderError(errors.New("Unable to send message"))
	}

	event := zone.SendMessage(user, text)
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
	zone, _ := chat.FindAvailableZone(geohash)
	if (zone.Geohash != geohash) {
		return c.RenderError(errors.New("Unable to join room"))
	}

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
