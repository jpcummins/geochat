package types

type User interface {
	ID() string
	Name() string
	NewConnection() (Connection, error)
	Disconnect(Connection) error
}
