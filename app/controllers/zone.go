package controllers

import (
    "golang.org/x/net/websocket"
    "github.com/revel/revel"
    "github.com/jpcummins/geochat/app/chat"
)

type Zone struct {
    *revel.Controller
}

func (c Zone) Message(user, text string, zone string) revel.Result {
    z, ok := chat.FindZone(zone)

    if (!ok) {
        return nil
    }

    event, err := z.SendMessage(user, text)

    if err != nil {
        return c.RenderError(err)
    }

    return c.RenderJson(event)
}

func (c Zone) Zone(user string) revel.Result {
    return c.Render(user)
}

func (c Zone) ZoneSocket(user string, ws *websocket.Conn) revel.Result {
    s := chat.SubscribeToZone("abc")
    defer s.Unsubscribe()

    for {
        select {
        case event := <-s.New:
            if websocket.JSON.Send(ws, &event) != nil {
                // They disconnected.
                return nil
            }
        }
    }
}