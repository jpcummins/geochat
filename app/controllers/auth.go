package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type AuthController struct {
	*revel.Controller
}

func (ac AuthController) Login(name string, lat float64, long float64) revel.Result {
	id, ok := ac.Session["subscription"]

	var subscription *chat.Subscription
	if ok {
		subscription, ok = (*chat.Subscribers).Get(id)

		if !ok {
			delete(ac.Session, "subscription")
		}
	}

	if !ok {
		user := &chat.User{Id: name, Name: name, Lat: lat, Long: long}
		subscription, err := chat.NewLocalSubscription(user)

		if err != nil {
			return ac.RenderError(err)
		}

		ac.Session["subscription"] = subscription.GetID()
	}

	return ac.RenderJson(subscription)
}
