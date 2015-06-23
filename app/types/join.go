package types

type ServerJoinJSON struct {
	Zone ServerJSON `json:"zone"`
	User ServerJSON `json:"user"`
}

type ClientJoinJSON struct {
	User ClientJSON `json:"user"`
}
