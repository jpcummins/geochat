package types

type Connection interface {
	Events() chan ClientEvent
	Ping()
}
