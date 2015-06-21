package types

type Event interface {
	ID() string
	Type() string
	WorldID() string
	World() World
	SetWorld(World)
	Data() EventData
}

type EventData interface {
	Type() string
	BeforePublish(Event) error
	OnReceive(Event) error
}
