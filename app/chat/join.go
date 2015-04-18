package chat

type Join struct {
	User *User `json:"user"`
}

func (j *Join) Type() string {
	return "join"
}
