package chat

import "errors"

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetUser(id string) (*User, error) {

	if id == "" {
		return nil, errors.New("Unauthorized")
	}

	return &User{Id: id, Name: id}, nil
}
