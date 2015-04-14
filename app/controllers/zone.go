package controllers

import (
    "golang.org/x/net/websocket"
    "github.com/revel/revel"
    "github.com/jpcummins/geochat/app/chatzone"
)

type Zone struct {
    *revel.Controller
}

func (c Zone) Message(user, text string) revel.Result {
    event := chatzone.Say(user, text)
    return c.RenderJson(event)
}

func (c Zone) Zone(user string) revel.Result {
    return c.Render(user)
}

func (c Zone) ZoneSocket(user string, ws *websocket.Conn) revel.Result {
    // Join the zone.
    subscription := chatzone.Subscribe()
    defer subscription.Cancel()

    chatzone.Join(user)
    defer chatzone.Leave(user)

    // Send down the archive.
    for _, event := range subscription.Archive {
        if websocket.JSON.Send(ws, &event) != nil {
            // They disconnected
            return nil
        }
    }

    // Now listen for new events from either the websocket or the chatzone.
    for {
        select {
        case event := <-subscription.New:
            if websocket.JSON.Send(ws, &event) != nil {
                // They disconnected.
                return nil
            }
        }
    }
    return nil
}