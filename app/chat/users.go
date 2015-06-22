package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Users struct {
	sync.RWMutex
	db      types.DB
	worldID string
	users   map[string]types.User
}

func newUsers(worldID string, db types.DB) *Users {
	return &Users{
		db:      db,
		worldID: worldID,
		users:   make(map[string]types.User),
	}
}

func (u *Users) User(id string) (types.User, error) {
	if cachedUser := u.FromCache(id); cachedUser != nil {
		return cachedUser, nil
	}

	return u.FromDB(id)
}

func (u *Users) FromCache(id string) types.User {
	u.RLock()
	defer u.RUnlock()
	return u.users[id]
}

func (u *Users) FromDB(id string) (types.User, error) {
	json, err := u.db.User(id, u.worldID)
	if err != nil {
		return nil, err
	}

	if json == nil {
		return nil, nil
	}

	user := u.FromCache(id)
	if user == nil {
		return nil, nil
	}

	user.Update(json)
	u.updateCache(user)
	return user, nil
}

func (u *Users) Save(user types.User) error {
	if err := u.db.SaveUser(user.ServerJSON()); err != nil {
		return err
	}

	u.updateCache(user)
	return nil
}

func (u *Users) updateCache(user types.User) {
	u.Lock()
	defer u.Unlock()
	u.users[user.ID()] = user
}
