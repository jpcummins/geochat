package types

type Events interface {
	New(string, ServerEventData) ServerEvent
}
