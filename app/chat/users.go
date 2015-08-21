package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Users struct {
	sync.RWMutex
	db    types.DB
	world types.World
	users map[string]types.User
}

func newUsers(world types.World, db types.DB) *Users {
	return &Users{
		db:    db,
		world: world,
		users: make(map[string]types.User),
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
	json, err := u.db.User(id, u.world.ID())
	if err != nil {
		return nil, err
	}

	if json == nil {
		return nil, nil
	}

	user := u.FromCache(id)
	if user == nil {
		user = newUser(id, u.world)
	}
	if err := user.Update(json); err != nil {
		return nil, err
	}

	u.UpdateCache(user)
	return user, nil
}

func (u *Users) Save(user types.User) error {
	json, ok := user.PubSubJSON().(*types.UserPubSubJSON)
	if !ok {
		return errors.New("Unable to serialize UserPubSubJSON.")
	}
	if err := u.db.SaveUser(json, u.world.ID()); err != nil {
		return err
	}

	u.UpdateCache(user)
	return nil
}

func (u *Users) UpdateCache(user types.User) {
	u.Lock()
	defer u.Unlock()
	u.users[user.ID()] = user
}
