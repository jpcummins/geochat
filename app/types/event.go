package types

type ServerEventType string

type ServerEvent interface {
	ID() string
	Type() ServerEventType
	World() World
	SetWorld(World)
	Data() ServerEventData
}

type ServerEventData interface {
	Type() ServerEventType
	BeforePublish(ServerEvent) error
	OnReceive(ServerEvent) error
}

type ClientEventType string

type ClientEvent interface {
	ID() string
	Type() ClientEventType
	Data() ClientEventData
}

type ClientEventData interface {
	Type() ClientEventType
	BeforeBroadcast(ClientEvent) error
}
