package chat

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"sync"
)

type Users struct {
	sync.RWMutex
	users map[string]*User
}

func NewUsers() *Users {
	u := &Users{
		users: make(map[string]*User),
	}
	return u
}

func (u *Users) Get(id string) (*User, bool) {
	if user, ok := u.cacheGet(id); ok {
		return user, ok
	}

	if user, ok := u.redisGet(id); ok {
		u.cacheSet(user)
		return user, ok
	}

	return nil, false
}

func (u *Users) Set(user *User) {
	u.cacheSet(user)
	u.redisSet(user)
}

func (u *Users) cacheGet(id string) (*User, bool) {
	u.RLock()
	user, found := u.users[id]
	u.RUnlock()
	return user, found
}

func (u *Users) cacheSet(user *User) {
	u.Lock()
	u.users[user.GetID()] = user
	u.Unlock()
}

func (u *Users) redisGet(id string) (*User, bool) {
	c := connection.Get()
	defer c.Close()

	usersJSON, err := redis.String(c.Do("HGET", "users", id))
	if err != nil {
		return nil, false
	}

	var user User
	if err := json.Unmarshal([]byte(usersJSON), &user); err != nil {
		return nil, false
	}
	return &user, true
}

func (u *Users) redisSet(user *User) {
	c := connection.Get()
	defer c.Close()
	eventJSON, err := json.Marshal(user)
	if err != nil {
		return
	}
	c.Do("HSET", "users", user.id, eventJSON)
}
