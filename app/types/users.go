package types

type Users interface {
	User(string) (User, error)
	UpdateUser(string) (User, error)
	SetUser(User) error
}
