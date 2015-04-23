package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type Auth struct {
	*revel.Controller
}

func (c Auth) Login(name string, lat float64, long float64) revel.Result {
	c.Session["user"] = name
	user := &chat.User{Id: name, Name: name, Lat: lat, Long: long}
	return c.RenderJson(user)
}
