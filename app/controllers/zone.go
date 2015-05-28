package controllers

import (
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

	subscription, found := (*chat.Subscribers).Get(subscriptionID)

	if !found {
		return zc.Redirect("/")
	}

	zc.subscription = subscription
	return nil
}

// Message action sends a message to those in the subscriber's zone.
func (zc *ZoneController) Message(text string) revel.Result {
	message := &chat.Message{Subscriber: zc.subscription, Text: text}
	event := chat.NewEvent(message)
	zone := zc.subscription.GetZone()
	zone.Publish(event)
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
	zc.subscription.Connect()
	zc.subscription.Events <- chat.NewEvent(zc.subscription.GetZone())
	closeConnection := make(chan bool)

	// Listen for client disconnects
	go func() {
		var msg string
		for {
			// The value of msg is ignored. Commands are not accepted over the websocket.
			if websocket.Message.Receive(ws, &msg) != nil {
				closeConnection <- true
				return
			}
		}
	}()

	// Send events to the WebSocket
	ping := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ping.C:
			zc.subscription.Events <- chat.NewEvent(&chat.Ping{})
		case event := <-zc.subscription.Events:
			if err := websocket.JSON.Send(ws, &event); err != nil {
				closeConnection <- true
			}
		case _ = <-closeConnection:
			zc.subscription.Disconnect()
			ws.Close()
			return nil
		}
	}
}
