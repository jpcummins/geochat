package types

type Users interface {
	FromCache(string) (User, error)
	FromDB(string) (User, error)
	Save(User) error
}
