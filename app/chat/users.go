package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Users struct {
	sync.RWMutex
	db    types.DB
	users map[string]types.User
}

func newUsers(db types.DB) *Users {
	return &Users{
		db:    db,
		users: make(map[string]types.User),
	}
}

func (u *Users) User(id string) (types.User, error) {
	user, found := u.localUser(id)
	if found {
		return user, nil
	}
	return u.UpdateUser(id)
}

func (u *Users) UpdateUser(id string) (types.User, error) {
	user := &User{}
	found, err := u.db.GetUser(id, user)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	u.localSetUser(user)
	return user, nil
}

func (u *Users) SetUser(user types.User) error {
	if err := u.db.SetUser(user); err != nil {
		return err
	}

	u.localSetUser(user)
	return nil
}

func (u *Users) localUser(id string) (types.User, bool) {
	u.RLock()
	defer u.RUnlock()
	user, found := u.users[id]
	return user, found
}

func (u *Users) localSetUser(user types.User) {
	u.Lock()
	defer u.Unlock()
	u.users[user.ID()] = user
}
