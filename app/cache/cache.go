package cache

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Cache struct {
	userMutex  sync.RWMutex
	zoneMutex  sync.RWMutex
	worldMutex sync.RWMutex
	users      map[string]types.User
	zones      map[string]types.Zone
	worlds     map[string]types.World
	db         types.DB
}

func NewCache(db types.DB) *Cache {
	return &Cache{
		users:  make(map[string]types.User),
		zones:  make(map[string]types.Zone),
		worlds: make(map[string]types.World),
		db:     db,
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

	if user != nil {
		err = c.localSetUser(user)
	}

	return user, err
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

	if zone != nil {
		err = c.localSetZone(zone)
	}

	return zone, err
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

func (c *Cache) World(id string) (types.World, error) {
	world, err := c.localWorld(id)
	if err != nil {
		return nil, err
	}

	if world != nil {
		return world, nil
	}

	world, err = c.db.GetWorld(id)
	if err != nil {
		return nil, err
	}

	if world != nil {
		err = c.localSetWorld(world)
	}

	return world, err
}

func (c *Cache) SetWorld(world types.World) error {
	if err := c.db.SetWorld(world); err != nil {
		return err
	}

	if err := c.localSetWorld(world); err != nil {
		return err
	}

	return nil
}

func (c *Cache) localWorld(id string) (types.World, error) {
	c.worldMutex.RLock()
	world := c.worlds[id]
	c.worldMutex.RUnlock()
	return world, nil
}

func (c *Cache) localSetWorld(world types.World) error {
	c.worldMutex.Lock()
	c.worlds[world.ID()] = world
	c.worldMutex.Unlock()
	return nil
}
