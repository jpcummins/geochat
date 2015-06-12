package types

type Cache interface {
	Get(id string) (User, error)
	Set(user User) error
}
