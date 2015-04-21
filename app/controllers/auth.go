package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
)

type Auth struct {
	*revel.Controller
}

func (c Auth) Login(name string, geohash string) revel.Result {
	c.Session["user"] = name
	user := &chat.User{Id: name, Name: name, Geohash: geohash}
	return c.RenderJson(user)
}
