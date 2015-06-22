package types

type Users interface {
	User(string) (User, error)
	Refresh(string) (User, error)
	Save(User) error
}
