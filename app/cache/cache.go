package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Cache struct {
	sync.RWMutex
	users map[string]types.User
	db    types.DB
}

func NewCache(db types.DB) *Cache {
	return &Cache{
		users: make(map[string]types.User),
		db:    db,
	}
}

func (c *Cache) User(id string) (types.User, error) {
	user, err := c.localUser(id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user, err = c.dbUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Cache) SetUser(user types.User) error {
	if err := c.localSetUser(user); err != nil {
		return err
	}

	if err := c.dbSetUser(user); err != nil {
		return err
	}

	return nil
}

func (c *Cache) localUser(id string) (types.User, error) {
	c.RLock()
	user := c.users[id]
	c.RUnlock()
	return user, nil
}

func (c *Cache) localSetUser(user types.User) error {
	c.Lock()
	c.users[user.ID()] = user
	c.Unlock()
	return nil
}

func (c *Cache) dbUser(id string) (types.User, error) {
	return c.db.GetUser(id)
}

func (c *Cache) dbSetUser(user types.User) error {
	return c.db.SetUser(user)
}

func (c *Cache) Zone(id string) (types.Zone, error) {
	return nil, nil
}

func (c *Cache) SetZone(zone types.Zone) error {
	return nil
}

func (c *Cache) GetZoneForUser(id string) (types.Zone, error) {
	return nil, nil
}
