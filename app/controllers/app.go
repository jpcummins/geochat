package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	if id, ok := c.Session["user_id"]; ok {

		user, err := chat.App.Users().User(id)
		if user != nil && err == nil {
			return c.Redirect("/chat")
		} else {
			delete(c.Session, "user_id")
		}
	}
	return c.Render()
}
