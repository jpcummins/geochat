package types

type Zone interface {
	ID() string
	Broadcast(Event)
	AddUser(User)
	RemoveUser(string)
}
