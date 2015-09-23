package controllers

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/jpcummins/geochat/app/types"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"io/ioutil"
	// "time"
)

// ZoneController is created for all requests handled by the Zone controller. It
// contains a handle to the chat package.
type ZoneController struct {
	*revel.Controller
	user  types.User
	world types.World
}

func init() {
	revel.InterceptMethod((*ZoneController).setSession, revel.BEFORE)
}

func (zc *ZoneController) setSession() revel.Result {
	userID, ok := zc.Session[userIDSessionKey]

	if !ok {
		zc.Response.Status = 401
		return zc.RenderError(errors.New("Unauthorized"))
	}

	user, err := chat.App.Users().User(userID)
	if err != nil {
		return zc.RenderError(err)
	}

	if user == nil {
		return zc.RenderError(errors.New("Invalid user."))
	}

	zc.user = user
	zc.world = chat.App
	return nil
}

// Message action sends a message to those in the subscriber's zone.
func (zc *ZoneController) Message() revel.Result {

	type messageData struct {
		Text string `json:"text"`
	}

	body, err := ioutil.ReadAll(zc.Request.Body)

	if err != nil {
		return zc.RenderError(err)
	}

	var data messageData
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return zc.RenderError(err)
	}

	if zc.user.ZoneID() == "" {
		zc.Redirect("/")
	}

	zone, err := zc.world.Zones().Zone(zc.user.ZoneID())
	if err != nil {
		panic(err)
	}

	if err := zone.Message(zc.user, data.Text); err != nil {
		panic(err)
	}

	return zc.RenderJson(nil)
}

// Command action is used to issue administrative commands
func (zc *ZoneController) Command(command string, args string) revel.Result {
	println(command, args)
	if err := zc.user.ExecuteCommand(command, args); err != nil {
		revel.ERROR.Printf("Error: %s\n", err.Error())
	}
	return nil
}

// Zone action renders the main chat interface
func (zc *ZoneController) Zone() revel.Result {
	userData, err := json.Marshal(zc.user.BroadcastJSON())
	userJSON := string(userData)

	if err != nil {
		return zc.RenderError(err)
	}

	return zc.Render(userJSON)
}

// ZoneSocket action handles WebSocket communication
func (zc *ZoneController) ZoneSocket(lat float64, long float64, ws *websocket.Conn) revel.Result {
	connection := zc.user.Connect()
	closeConnection := make(chan bool)

	if zc.user.ZoneID() == "" {
		if _, err := zc.world.Join(zc.user); err != nil {
			panic(err)
		}
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
	// ping := time.NewTicker(30 * time.Second)
	events := connection.Events()

	for {
		select {
		// case <-ping.C:
		// 	connection.Ping()
		case event := <-events:
			if err := websocket.JSON.Send(ws, &event); err != nil {
				closeConnection <- true
			}
		case _ = <-closeConnection:
			// Leave the chat room
			zoneID := zc.user.ZoneID()

			if zoneID != "" {
				zone, err := zc.world.Zones().Zone(zoneID)
				if err != nil {
					panic(err)
				}
				if err := zone.Leave(zc.user); err != nil {
					panic(err)
				}
			}

			zc.user.Disconnect(connection)
			ws.Close()
			return nil
		}
	}
}
