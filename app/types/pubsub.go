package types

type PubSub interface {
	Publish(ServerEvent) error
	Subscribe() <-chan ServerEvent
}
