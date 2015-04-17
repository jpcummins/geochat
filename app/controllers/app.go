package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	User *chat.User
}

func (c App) Index(user string) revel.Result {
	c.Validation.Required(user)
	return c.Redirect("/zone/zone?user=%s", user)
}
