package types

type Event interface {
	ID() string
	Type() string
	Zone() Zone
	Data() EventData
}

type EventData interface {
	Type() string
	BeforePublish(Event) error
	OnReceive(Event) error
}
