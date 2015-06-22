package types

type Events interface {
	NewServerEvent(string, ServerEventData) ServerEvent
	NewClientEvent(string, ClientEventData) ClientEvent
}
