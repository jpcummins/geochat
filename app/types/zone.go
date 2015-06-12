package types

type Zone interface {
	ID() string
	Count() int
	IsOpen() bool
	Broadcast(Event)
	AddUser(User) error
	RemoveUser(string) error
}
