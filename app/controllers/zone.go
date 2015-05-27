package controllers

import (
	"fmt"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"time"
)

// ZoneController is created for all requests handled by the Zone controller. It
// contains a handle to the chat package.
type ZoneController struct {
	*revel.Controller
	subscription *chat.Subscription
}

func init() {
	revel.InterceptMethod((*ZoneController).setSession, revel.BEFORE)
}

func (zc *ZoneController) setSession() revel.Result {
	subscriptionID, ok := zc.Session["subscription"]

	if !ok {
		zc.Redirect("/")
	}

	subscription := chat.GetSubscription(subscriptionID)

	if subscription == nil {
		return zc.Redirect("/")
	}

	zc.subscription = subscription
	return nil
}

// Message action sends a message to those in the subscriber's zone.
func (zc *ZoneController) Message(text string) revel.Result {
	message := &chat.Message{User: zc.subscription.GetUser(), Text: text}
	event := chat.NewEvent(message)
	zc.subscription.GetZone().Publish(event)
	return zc.RenderJson(event)
}

// Command action is used to issue administrative commands
func (zc *ZoneController) Command(command string) revel.Result {
	json, err := zc.subscription.ExecuteCommand(command)
	if err != nil {
		return zc.RenderError(err)
	}
	return zc.RenderJson(json)
}

// Zone action renders the main chat interface
func (zc *ZoneController) Zone() revel.Result {
	return zc.Render()
}

// ZoneSocket action handles WebSocket communication
func (zc *ZoneController) ZoneSocket(ws *websocket.Conn) revel.Result {
	zc.subscription.SetOnline()

	// Listen for client disconnects
	go func() {
		var msg string
		for {
			if websocket.Message.Receive(ws, &msg) != nil {
				zc.subscription.SetOffline()
				ws.Close()
				return
			}
		}
	}()

	zone := zc.subscription.GetZone()

	// Find a zone if one hasn't been set
	if zone == nil {
		user := zc.subscription.GetUser()
		var err error
		zone, err = chat.GetOrCreateAvailableZone(user.Lat, user.Long)
		if err != nil {
			return zc.RenderError(err)
		}
		zc.subscription.SetZone(zone)
	}

	zc.subscription.Events <- chat.NewEvent(zone)
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			zc.subscription.Events <- chat.NewEvent(&chat.Ping{})
		case event := <-zc.subscription.Events:
			if err := websocket.JSON.Send(ws, &event); err != nil {
				zc.subscription.SetOffline()
				ws.Close()
				return nil
			}
		}
	}
}
