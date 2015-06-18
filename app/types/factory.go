package types

type EventFactory interface {
	New(string, EventData) (Event, error)
}
