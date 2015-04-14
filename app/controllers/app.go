package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index(user string) revel.Result {
    c.Validation.Required(user)
    return c.Redirect("/zone/zone?user=%s", user)
}