package types

type EventJSON interface {
	ClientJSON() ClientJSON
	ServerJSON() ServerJSON
	Update(ServerJSON) error
}

type ClientJSON interface{}

type ServerJSON interface {
	Key() string
	WorldKey() string
}

// Base implementations

type BaseClientJSON struct {
	ClientJSON `json:"-"`
	ID         string `json:"id"`
}

type BaseServerJSON struct {
	ServerJSON `json:"-"`
	ID         string `json:"id"`
	WorldID    string `json:"-"`
}

func (json *BaseServerJSON) WorldKey() string {
	return json.WorldID
}

func (json *BaseServerJSON) Key() string {
	return json.ID
}
