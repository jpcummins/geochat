package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"github.com/TomiHiltunen/geohash-golang"
)

type Zone struct {
	*revel.Controller
}

func (c Zone) Lookup(lat float64, long float64) revel.Result {
	return c.RenderText(geohash.EncodeWithPrecision(lat, long, 3))
}

func (c Zone) Message(zone string, text string) revel.Result {

    user, err := chat.GetUser(c.Session["user"])
    if err != nil {
        return c.RenderError(err)
    }

	z, ok := chat.FindZone(zone)

	if !ok {
		return nil
	}

	event, err := z.SendMessage(user, text)

	if err != nil {
		return c.RenderError(err)
	}

	return c.RenderJson(event)
}

func (c Zone) Zone(zone string) revel.Result {
	box := geohash.Decode(zone)
	return c.Render(zone, box)
}

func (c Zone) ZoneSocket(zone string, ws *websocket.Conn) revel.Result {
	s := chat.SubscribeToZone(zone)
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
