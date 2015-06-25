package types

type Connection interface {
	Events() <-chan BroadcastEvent
	Ping()
}
