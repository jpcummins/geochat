package controllers

import (
	"errors"
	"fmt"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"time"
)

type Zone struct {
	*revel.Controller
	chat.Chat
}

func (c Zone) Subscribe(lat float64, long float64) revel.Result {
	user, err := chat.GetUser(c.Session["user"])
	if err != nil {
		return c.RenderError(err)
	}
	zone, err := chat.GetOrCreateAvailableZone(lat, long)
	if err != nil {
		return c.RenderError(err)
	}

	subscription := c.Subscribers.Add(user, zone)
	return c.RenderJson(subscription)
}

func (c Zone) Message(geochat *chat.Chat, subscriptionId string, text string) revel.Result {
	subscription := c.Subscribers.Get(subscriptionId)
	if subscription == nil {
		return c.RenderError(errors.New("Invalid subscription"))
	}

	event := subscription.Broadcast(text)
	return c.RenderJson(event)
}

func (c Zone) Command(geochat *chat.Chat, subscriptionId string, command string) revel.Result {
	subscription := c.Subscribers.Get(subscriptionId)
	if subscription == nil {
		return c.RenderError(errors.New("Invalid subscription"))
	}

	json, err := chat.ExecuteCommand(command, subscription)
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJson(json)
}

func (c Zone) Zone(subscriptionId string) revel.Result {
	subscription := c.Subscribers.Get(subscriptionId)
	if subscription == nil {
		return c.Redirect("/")
	}
	return c.Render(subscriptionId)
}

func (c Zone) ZoneSocket(subscriptionId string, ws *websocket.Conn) revel.Result {
	subscription := c.Subscribers.Get(subscriptionId)
	if subscription == nil {
		return c.RenderError(errors.New("Invalid subscription"))
	}

	subscription.Activate()

	// Listen for client disconnects
	go func() {
		var msg string
		for {
			if websocket.Message.Receive(ws, &msg) != nil {
				println("close1")
				subscription.Deactivate()
				ws.Close()
				return
			}
		}
	}()

	zone, err := chat.GetOrCreateZone(subscription.Zonehash)
	if err != nil {
		return c.RenderError(err)
	}

	zoneInfo := &chat.ZoneInfo{
		ID:          zone.Zonehash,
		Boundary:    zone.Boundary,
		Archive:     zone.GetArchive(10),
		Subscribers: zone.GetSubscribers(),
	}

	subscription.Events <- chat.NewEvent(zoneInfo)

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			subscription.Events <- chat.NewEvent(&chat.Ping{})
		case event := <-subscription.Events:

			if !subscription.IsOnline {
				println("close2")
				ws.Close()
				return nil
			}

			fmt.Printf("%+v\n", event)
			if err := websocket.JSON.Send(ws, &event); err != nil {
				println("close3", err.Error())
				subscription.Deactivate()
				ws.Close()
				return nil
			}
		}
	}
}
