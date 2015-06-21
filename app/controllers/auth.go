package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type AuthController struct {
	*revel.Controller
}

var userIDSessionKey = "user_id"

func (ac AuthController) Login(name string, lat float64, long float64) revel.Result {
	id, ok := ac.Session[userIDSessionKey]

	user, err := chat.App.Users().User(id)
	if err != nil {
		delete(ac.Session, userIDSessionKey)
	}

	if user == nil {
		user = chat.NewUser(name, name, NewLatLng(lat, long))
		if err := chat.App.Users().SetUser(user); err != nil {
			delete(ac.Session, userIDSessionKey)
			return ac.RenderError(err)
		}
		ac.Session[userIDSessionKey] = user.GetID()
	}

	return ac.RenderJson(user)
}
