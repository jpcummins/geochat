package chat

type Join struct {
	User string `json:"user"`
}

func (j *Join) Type() string {
	return "join"
}
