package types

type PubSub interface {
	Publish(PubSubEvent) error
	Subscribe() <-chan PubSubEvent
}

type PubSubEventType string

type PubSubEvent interface {
	ID() string
	Type() PubSubEventType
	Data() PubSubEventData

	World() World
	SetWorld(World)
}

type PubSubEventData interface {
	Type() PubSubEventType
	BeforePublish(PubSubEvent) error
	OnReceive(PubSubEvent) error
}

type PubSubDataType string

type PubSubSerializable interface {
	PubSubJSON() PubSubJSON
	Update(PubSubJSON) error
}

type PubSubJSON interface {
	Type() PubSubDataType
}
