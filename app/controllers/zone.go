package controllers

import (
	"errors"
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"time"
)

type Zone struct {
	*revel.Controller
}

func (c Zone) Subscribe(lat float64, long float64) revel.Result {
	user, err := chat.GetUser(c.Session["user"])
	if err != nil {
		return c.RenderError(err)
	}

	geohash := gh.EncodeWithPrecision(lat, long, 6)

	zone, err := chat.FindAvailableZone(geohash)
	if err != nil {
		return c.RenderError(err)
	}

	subscription := zone.Subscribe(user)
	return c.RenderJson(subscription)
}

func (c Zone) Message(subscriptionId string, text string) revel.Result {
	user, err := chat.GetUser(c.Session["user"])
	if err != nil {
		return c.RenderError(err)
	}

	subscription := chat.GetSubscription(subscriptionId)
	event := subscription.Zone.SendMessage(user, text)
	return c.RenderJson(event)
}

func (c Zone) Command(subscriptionId string, command string) revel.Result {
	subscription := chat.GetSubscription(subscriptionId)
	json, err := chat.ExecuteCommand(command, subscription)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(json)
}

func (c Zone) Zone(subscriptionId string) revel.Result {
	subscription := chat.GetSubscription(subscriptionId)
	if subscription == nil || subscription.Zone == nil {
		c.Redirect("/")
	}

	zonehash := subscription.Zone.GetZonehash()
	box := gh.Decode(subscription.Zone.Geohash) // TODO: incorporate subhash
	return c.Render(zonehash, box, subscriptionId)
}

func (c Zone) ZoneSocket(subscriptionId string, ws *websocket.Conn) revel.Result {
	subscription := chat.GetSubscription(subscriptionId)
	if subscription == nil {
		c.RenderError(errors.New("Invalid subscription"))
	}

	// Listen for client disconnects
	go func() {
		var msg string
		for {
			if websocket.Message.Receive(ws, &msg) != nil {
				subscription.Zone.Unsubscribe(subscription)
				return
			}
		}
	}()

	// Send zone information
	subscription.Events <- chat.NewEvent(subscription.Zone)

	// Send the archive
	archive, _ := subscription.Zone.GetArchive(10)
	subscription.Events <- chat.NewEvent(archive)

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			subscription.Events <- &chat.Event{Type: "ping", Data: nil, Timestamp: int(time.Now().Unix())}
		case event := <-subscription.Events:
			if websocket.JSON.Send(ws, &event) != nil {
				subscription.Zone.Unsubscribe(subscription)
				return nil
			}
		}
	}
}
