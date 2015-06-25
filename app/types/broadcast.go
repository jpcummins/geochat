package types

type BroadcastEventType string

type BroadcastEvent interface {
	ID() string
	Type() BroadcastEventType
	Data() BroadcastEventData
}

type BroadcastEventData interface {
	Type() BroadcastEventType
	BeforeBroadcastToUser(User, BroadcastEvent) (bool, error)
}

type BroadcastSerializable interface {
	BroadcastJSON() interface{}
}
