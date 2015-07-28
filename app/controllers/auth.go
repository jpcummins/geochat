package controllers

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/jpcummins/geochat/app/types"
	"github.com/revel/revel"
	"net/http"
	"net/url"
)

type AuthController struct {
	*revel.Controller
}

var userIDSessionKey = "user_id"

func (ac AuthController) Login(fbID string, lat float64, long float64, authToken string) revel.Result {
	id, ok := ac.Session[userIDSessionKey]

	var user types.User
	var err error

	if fbID == "" || authToken == "" {
		return ac.RenderError(errors.New("Invalid Facebook credentials."))
	}

	if ok {
		user, err = chat.App.Users().User(id)
		if err != nil {
			delete(ac.Session, userIDSessionKey)
		}
	}

	fbUser := map[string]interface{}{}
	fbUserResponse, _ := http.Get("https://graph.facebook.com/me?access_token=" + url.QueryEscape(authToken))
	defer fbUserResponse.Body.Close()
	if err := json.NewDecoder(fbUserResponse.Body).Decode(&fbUser); err != nil {
		revel.ERROR.Println(err)
		return ac.RenderError(err)
	}

	fbPicture := map[string]interface{}{}
	fbPictureResponse, _ := http.Get("https://graph.facebook.com/v2.4/" + fbID + "/picture?type=large&redirect=false&access_token=" + url.QueryEscape(authToken))
	defer fbPictureResponse.Body.Close()
	if err := json.NewDecoder(fbPictureResponse.Body).Decode(&fbPicture); err != nil {
		revel.ERROR.Println(err)
		return ac.RenderError(err)
	}

	fbPictureData := fbPicture["data"].(map[string]interface{})
	fbPictureURL := fbPictureData["url"].(string)
	name := fbUser["name"].(string)
	firstName := fbUser["first_name"].(string)
	lastName := fbUser["last_name"].(string)
	timezone := fbUser["timezone"].(float64)
	email := fbUser["email"].(string)

	if isAnyNil(fbPictureData, fbPictureURL, name, firstName, lastName, timezone, email) {
		return ac.RenderError(errors.New("Unable to parse FB data"))
	}

	if user == nil {
		user, err = chat.App.NewUser()
		if err != nil {
			delete(ac.Session, userIDSessionKey)
			ac.RenderError(err)
		}
		ac.Session[userIDSessionKey] = user.ID()
	}

	user.SetName(name)
	user.SetFirstName(firstName)
	user.SetLastName(lastName)
	user.SetTimezone(timezone)
	user.SetEmail(email)
	user.SetFBID(fbID)
	user.SetFBAccessToken(authToken)
	user.SetFBPictureURL(fbPictureURL)
	user.SetLocation(lat, long)

	return ac.RenderJson(user)
}

func (ac AuthController) Logout() revel.Result {
	delete(ac.Session, userIDSessionKey)
	return ac.Redirect("/")
}

func isAnyNil(args ...interface{}) bool {
	for _, arg := range args {
		if arg == nil {
			return true
		}
	}
	return false
}
