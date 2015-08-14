package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jpcummins/geochat/app/chat"
	"github.com/jpcummins/geochat/app/types"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"strconv"
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
		err := errors.New("Invalid Facebook credentials.")
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}

	if authToken == "test" {
		if user, err = chat.App.Users().User(fbID); err == nil && user != nil {
			ac.Session[userIDSessionKey] = user.ID()
			return ac.RenderJson(user.PubSubJSON())
		}
		return ac.RenderError(errors.New("Unable to find test user"))
	}

	if ok {
		user, err = chat.App.Users().User(id)
		if err != nil {
			delete(ac.Session, userIDSessionKey)
			revel.ERROR.Printf("Error retrieving user: %s\n", err.Error())
			return ac.RenderError(err)
		}
	}

	appSecret, found := revel.Config.String("app.fbAppSecret")
	if !found {
		err := errors.New("Unable to get app secret")
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}

	key := []byte(appSecret)
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(authToken))
	appSecretProof := hex.EncodeToString(sig.Sum(nil))

	// Get user's info
	url := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&appsecret_proof=%s", authToken, appSecretProof)
	fbUserResponse, err := http.Get(url)
	if err != nil {
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}
	if fbUserResponse.StatusCode != 200 {
		err := errors.New("FB returned a non-200 status for /me:" + strconv.Itoa(fbUserResponse.StatusCode))
		revel.ERROR.Printf("Error: %s\n", err.Error())
		contents, _ := ioutil.ReadAll(fbUserResponse.Body)
		revel.ERROR.Printf("Body: " + string(contents))
		return ac.RenderError(err)
	}

	defer fbUserResponse.Body.Close()
	fbUser := map[string]interface{}{}
	if err := json.NewDecoder(fbUserResponse.Body).Decode(&fbUser); err != nil {
		revel.ERROR.Println(err)
		return ac.RenderError(err)
	}

	// Get user's profile pic
	url = fmt.Sprintf("https://graph.facebook.com/v2.4/%s/picture?type=large&redirect=false&access_token=%s&appsecret_proof=%s", fbID, authToken, appSecretProof)
	fbPictureResponse, err := http.Get(url)
	if err != nil {
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}
	if fbPictureResponse.StatusCode != 200 {
		err := errors.New("FB returned a non-200 status for /v2.4/[id]/picture:" + string(fbPictureResponse.StatusCode))
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}

	fbPicture := map[string]interface{}{}
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
		err := errors.New("Unable to parse FB data")
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}

	if user == nil {
		user, err = chat.App.NewUser(fbID)
		if err != nil {
			delete(ac.Session, userIDSessionKey)
			revel.ERROR.Printf("Error: %s\n", err.Error())
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
	chat.App.Users().Save(user)

	return ac.RenderJson(user.PubSubJSON())
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
