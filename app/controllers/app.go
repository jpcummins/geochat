package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	User *chat.User
}

func (c App) Index() revel.Result {
	user, _ := chat.GetUser(c.Session["user"])
	return c.Render(user)
}
