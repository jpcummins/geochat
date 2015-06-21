package types

type Connection interface {
	Events() chan Event
	Ping()
}
