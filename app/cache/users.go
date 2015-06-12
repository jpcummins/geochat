package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type UserCache struct {
	sync.RWMutex
	users map[string]types.User
	db    types.DB
}

func NewUserCache(db types.DB) *UserCache {
	return &UserCache{
		users: make(map[string]types.User),
		db:    db,
	}
}

func (uc *UserCache) Get(id string) (types.User, error) {
	user, err := uc.localGet(id)
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

func (uc *UserCache) Set(user types.User) error {
	if err := uc.localSet(user); err != nil {
		return err
	}

	if err := uc.dbSet(user); err != nil {
		return err
	}

	return nil
}

func (uc *UserCache) localGet(id string) (types.User, error) {
	uc.RLock()
	user := uc.users[id]
	uc.RUnlock()
	return user, nil
}

func (uc *UserCache) localSet(user types.User) error {
	uc.Lock()
	uc.users[user.ID()] = user
	uc.Unlock()
	return nil
}

func (uc *UserCache) dbGet(id string) (types.User, error) {
	return uc.db.GetUser(id)
}

func (uc *UserCache) dbSet(user types.User) error {
	return uc.db.SetUser(user)
}
