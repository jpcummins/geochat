package controllers

import (
	"github.com/jpcummins/geochat/app/chat"
	"github.com/revel/revel"
	"net/http"
)

var AuthorizedFilter = func(c *revel.Controller, fc []revel.Filter) {
	user := c.Params.Values.Get("user")

	if user == "" {
		c.Response.Status = http.StatusUnauthorized

		// TODO: Return JSON instead
		c.RenderText("unauthorized")
	}

	if zc, ok := c.AppController.(*Zone); ok {
		zc.User = &chat.User{Name: user}
	}

	// Execute the next filter stage
	fc[0](c, fc[1:])
}
