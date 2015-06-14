package types

type Zone interface {
	ID() string
	Count() int
	IsOpen() bool
	Broadcast(Event)
	AddUser(User)
	RemoveUser(string)
}
