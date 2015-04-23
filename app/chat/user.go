package chat

import "errors"

type User struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Lat   float64 `json:"-"`
	Long  float64 `json:"-"`
	IsBot bool    `json:"is_bot,omitempty"`
}

func GetUser(id string) (*User, error) {

	if id == "" {
		return nil, errors.New("Unauthorized")
	}

	return &User{Id: id, Name: id}, nil
}
