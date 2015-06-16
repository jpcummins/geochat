package types

type PubSub interface {
	Publish(Event) error
	Subscribe() <-chan Event
}
