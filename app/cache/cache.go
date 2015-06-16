package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Cache struct {
	userMutex sync.RWMutex
	zoneMutex sync.RWMutex
	users     map[string]types.User
	zones     map[string]types.Zone
	db        types.DB
}

func NewCache(db types.DB) *Cache {
	return &Cache{
		users: make(map[string]types.User),
		zones: make(map[string]types.Zone),
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

	user, err = c.db.GetUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Cache) SetUser(user types.User) error {
	if err := c.db.SetUser(user); err != nil {
		return err
	}

	if err := c.localSetUser(user); err != nil {
		return err
	}

	return nil
}

func (c *Cache) localUser(id string) (types.User, error) {
	c.userMutex.RLock()
	user := c.users[id]
	c.userMutex.RUnlock()
	return user, nil
}

func (c *Cache) localSetUser(user types.User) error {
	c.userMutex.Lock()
	c.users[user.ID()] = user
	c.userMutex.Unlock()
	return nil
}

func (c *Cache) Zone(id string) (types.Zone, error) {
	zone, err := c.localZone(id)
	if err != nil {
		return nil, err
	}

	if zone != nil {
		return zone, nil
	}

	zone, err = c.db.GetZone(id)
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (c *Cache) SetZone(zone types.Zone) error {
	if err := c.db.SetZone(zone); err != nil {
		return err
	}

	if err := c.localSetZone(zone); err != nil {
		return err
	}

	return nil
}

func (c *Cache) localZone(id string) (types.Zone, error) {
	c.zoneMutex.RLock()
	zone := c.zones[id]
	c.zoneMutex.RUnlock()
	return zone, nil
}

func (c *Cache) localSetZone(zone types.Zone) error {
	c.zoneMutex.Lock()
	c.zones[zone.ID()] = zone
	c.zoneMutex.Unlock()
	return nil
}

func (c *Cache) GetZoneForUser(id string) (types.Zone, error) {
	return nil, nil
}
