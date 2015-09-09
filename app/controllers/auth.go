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
	"math/rand"
	"net/http"
	"strconv"
)

type AuthController struct {
	*revel.Controller
}

var userIDSessionKey = "user_id"

func (ac AuthController) LoginJson() revel.Result {

	type authData struct {
		FBID      string  `json:"fbID"`
		Lat       float64 `json:"lat"`
		Long      float64 `json:"long"`
		AuthToken string  `json:"authToken"`
	}

	body, err := ioutil.ReadAll(ac.Request.Body)

	if err != nil {
		return ac.RenderError(err)
	}

	var data authData
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return ac.RenderError(err)
	}

	return ac.Login(data.FBID, data.Lat, data.Long, data.AuthToken)
}

func (ac AuthController) Login(fbID string, lat float64, long float64, authToken string) revel.Result {
	revel.INFO.Printf("Login from %s. Lat: %f Lng: %f, auth: %s", fbID, lat, long, authToken)

	var user types.User
	var err error

	if authToken == "test" {
		user, err := chat.App.NewUser(strconv.Itoa(rand.Intn(10000000)))
		if err != nil {
			return ac.RenderError(err)
		}

		user.SetFirstName("Test")
		user.SetLastName("User")
		user.SetName("Test User")
		user.SetLocation(lat, long)

		url := fmt.Sprintf("http://uifaces.com/api/v1/random")
		avatarResponse, err := http.Get(url)
		if err != nil {
			revel.ERROR.Printf("Error: %s\n", err.Error())
			return ac.RenderError(err)
		}

		avatarJs := map[string]interface{}{}
		if err := json.NewDecoder(avatarResponse.Body).Decode(&avatarJs); err != nil {
			revel.ERROR.Println(err)
			return ac.RenderError(err)
		}
		avatarSizes := avatarJs["image_urls"].(map[string]interface{})
		user.SetFBPictureURL(avatarSizes["normal"].(string))

		chat.App.Users().Save(user)
		ac.Session[userIDSessionKey] = user.ID()
		return ac.RenderJson(user.PubSubJSON())
	}

	if fbID == "" || authToken == "" {
		err := errors.New("Invalid Facebook credentials.")
		revel.ERROR.Printf("Error: %s\n", err.Error())
		return ac.RenderError(err)
	}

	user, err = chat.App.Users().User(fbID)
	if err != nil {
		delete(ac.Session, userIDSessionKey)
		revel.ERROR.Printf("Error retrieving user: %s\n", err.Error())
		return ac.RenderError(err)
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
	ac.Session[userIDSessionKey] = user.ID()
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
