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
		user, err := chat.NewUser(lat, long, name)

		if err != nil {
			return ac.RenderError(err)
		}

		ac.Session["user_id"] = user.GetID()
	}

	return ac.RenderJson(user)
}
