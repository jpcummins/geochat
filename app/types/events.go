package types

type Events interface {
	New(string, EventData) (Event, error)
}
