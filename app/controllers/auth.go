package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type AuthController struct {
	*revel.Controller
}

func (ac AuthController) Login(name string, lat float64, long float64) revel.Result {
	id, ok := ac.Session["user_id"]

	var user *chat.User
	if ok {
		user, ok = (*chat.UserCache).Get(id)

		if !ok {
			delete(ac.Session, "user_id")
		}
	}

	if !ok {
		println("creating new user")
		user := chat.NewUser(lat, long, name)
		if _, err := user.JoinNextAvailableZone(); err != nil {
			return ac.RenderError(err)
		}

		ac.Session["user_id"] = user.GetID()
	}

	return ac.RenderJson(user)
}
