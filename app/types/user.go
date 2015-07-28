package types

type User interface {
	PubSubSerializable
	BroadcastSerializable

	ID() string

	Name() string
	SetName(string)

	FirstName() string
	SetFirstName(string)

	LastName() string
	SetLastName(string)

	Timezone() float64
	SetTimezone(float64)

	Email() string
	SetEmail(string)

	Location() LatLng
	SetLocation(float64, float64)

	ZoneID() string
	SetZoneID(string)

	FBID() string
	SetFBID(string)

	FBAccessToken() string
	SetFBAccessToken(string)

	FBPictureURL() string
	SetFBPictureURL(string)

	Broadcast(BroadcastEventData)
	ExecuteCommand(command string, args string) error

	Connect() Connection
	Disconnect(Connection)
}

type UserPubSubJSON struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	FirstName     string      `json:"first_name"`
	LastName      string      `json:"last_name"`
	Timezone      float64     `json:"timezone"`
	Email         string      `json:"email"`
	Location      *LatLngJSON `json:"location"`
	ZoneID        string      `json:"zone_id"`
	FBID          string      `json:"fb_id"`
	FBAccessToken string      `json:"fb_access_token"`
	FBPictureURL  string      `json:"fb_picture_url"`
}

func (psUser *UserPubSubJSON) Type() PubSubDataType {
	return "user"
}

type UserBroadcastJSON struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Location     *LatLngJSON `json:"location"`
	FBPictureURL string      `json:"fb_picture_url"`
}
