package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/jpcummins/geochat/app/types"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"time"
)

// ZoneController is created for all requests handled by the Zone controller. It
// contains a handle to the chat package.
type ZoneController struct {
	*revel.Controller
	user types.User
}

func init() {
	revel.InterceptMethod((*ZoneController).setSession, revel.BEFORE)
}

func (zc *ZoneController) setSession() revel.Result {
	userID, ok := zc.Session[userIDSessionKey]

	if !ok {
		zc.Redirect("/")
	}

	user, err := chat.App.Users().User(userID)
	if err != nil {
		return zc.RenderError(err)
	}

	if user == nil {
		return zc.Redirect("/")
	}

	zc.user = user
	return nil
}

// Message action sends a message to those in the subscriber's zone.
func (zc *ZoneController) Message(text string) revel.Result {

	if zc.user.Zone() == nil {
		zc.Redirect("/")
	}

	event, err := zc.user.Zone().Message(zc.user, text)

	if err != nil {
		return zc.RenderError(err)
	}

	return zc.RenderJson(event)
}

// Command action is used to issue administrative commands
func (zc *ZoneController) Command(command string) revel.Result {
	return zc.Redirect("/")
}

// Zone action renders the main chat interface
func (zc *ZoneController) Zone() revel.Result {
	return zc.Render()
}

// ZoneSocket action handles WebSocket communication
func (zc *ZoneController) ZoneSocket(ws *websocket.Conn) revel.Result {
	var zone types.Zone
	var err error

	connection := zc.user.Connect()
	closeConnection := make(chan bool)

	if zc.user.Zone() == nil {
		zone, err = chat.App.FindOpenZone(zc.user)
		if err != nil {
			println(err.Error())
			return zc.RenderError(err)
		}
		if _, err := zone.Join(zc.user); err != nil {
			println(err.Error())
			return zc.RenderError(err)
		}
	} else {
		zone = zc.user.Zone()
	}

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
	events := connection.Events()

	for {
		select {
		case <-ping.C:
			connection.Ping()
		case event := <-events:
			if err := websocket.JSON.Send(ws, &event); err != nil {
				closeConnection <- true
			}
		case _ = <-closeConnection:
			// Leave the chat room
			zone := zc.user.Zone()
			if zone != nil {
				if _, err := zone.Leave(zc.user); err != nil {
					println(err.Error())
				}
			}

			zc.user.Disconnect(connection)
			ws.Close()
			return nil
		}
	}
}
