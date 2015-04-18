package controllers

import (
	"github.com/revel/revel"
	"github.com/jpcummins/geochat/app/chat"
)

type Auth struct {
	*revel.Controller
}

func (c Auth) Login(name string) revel.Result {
	c.Session["user"] = name
	user := &chat.User{Name: name}
	return c.RenderJson(user)
}