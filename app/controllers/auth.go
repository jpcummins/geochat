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
		subscription = chat.GetSubscription(id)
		if subscription == nil {
			delete(ac.Session, "subscription")
			ok = !ok
		}
	}

	if !ok {
		user := &chat.User{Id: name, Name: name, Lat: lat, Long: long}
		subscription = chat.NewSubscription(user)
		ac.Session["subscription"] = subscription.GetID()
	}

	return ac.RenderJson(subscription)
}
