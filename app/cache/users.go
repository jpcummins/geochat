package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type UserCache struct {
	sync.RWMutex
	users      map[string]types.User
	connection types.Connection
}

func NewUserCache(connection types.Connection) *UserCache {
	return &UserCache{
		users:      make(map[string]types.User),
		connection: connection,
	}
}

func (uc *UserCache) Get(id string) (types.User, error) {
	user, err := uc.LocalGet(id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user, err = uc.dbGet(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserCache) Set(user types.User) {
	uc.LocalSet(user)
	uc.dbSet(user)
}

func (uc *UserCache) LocalGet(id string) (types.User, error) {
	uc.RLock()
	user := uc.users[id]
	uc.RUnlock()
	return user, nil
}

func (uc *UserCache) LocalSet(user types.User) {
	uc.Lock()
	uc.users[user.ID()] = user
	uc.Unlock()
}

func (uc *UserCache) dbGet(id string) (types.User, error) {
	return uc.connection.GetUser(id)
}

func (uc *UserCache) dbSet(user types.User) error {
	return uc.connection.SetUser(user)
}
