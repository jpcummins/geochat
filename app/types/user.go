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

	Locality() string
	SetLocality(string)

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
	FirstName     string      `json:"firstName"`
	LastName      string      `json:"lastName"`
	Timezone      float64     `json:"timezone"`
	Email         string      `json:"email"`
	Location      *LatLngJSON `json:"location"`
	Locality      string      `json:"locality"`
	ZoneID        string      `json:"zoneId"`
	FBID          string      `json:"fbID"`
	FBAccessToken string      `json:"fbAccessToken"`
	FBPictureURL  string      `json:"fbPictureUrl"`
}

func (psUser *UserPubSubJSON) Type() PubSubDataType {
	return "user"
}

type UserBroadcastJSON struct {
	ID           string      `json:"id"`
	FirstName    string      `json:"firstName"`
	Locality     string      `json:"locality"`
	Location     *LatLngJSON `json:"location"`
	FBPictureURL string      `json:"fbPictureUrl"`
}
