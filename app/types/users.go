package types

type Users interface {
	User(string) (User, error)
	FromCache(string) User
	UpdateCache(User)
	FromDB(string) (User, error)
	Save(User) error
}
