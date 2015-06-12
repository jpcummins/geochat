package types

type User interface {
	ID() string
	Name() string
	Zone() Zone
	NewConnection() (Connection, error)
	Disconnect(Connection) error
}
