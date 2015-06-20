package types

type Chat interface {
	DB() DB
	PubSub() PubSub
	Events() EventFactory
	World() World
}
