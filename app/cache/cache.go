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
	zones      map[string](map[string]types.Zone)
	worlds     map[string]types.World
	db         types.DB
}

func NewCache(db types.DB) *Cache {
	return &Cache{
		users:  make(map[string]types.User),
		zones:  make(map[string](map[string]types.Zone)),
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

	return c.UpdateUser(id)
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

func (c *Cache) UpdateUser(id string) (types.User, error) {
	user, err := c.db.GetUser(id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = c.localSetUser(user)
	}

	return user, err
}

func (c *Cache) localUser(id string) (types.User, error) {
	c.userMutex.RLock()
	defer c.userMutex.RUnlock()
	return c.users[id], nil
}

func (c *Cache) localSetUser(user types.User) error {
	c.userMutex.Lock()
	defer c.userMutex.Unlock()
	c.users[user.ID()] = user
	return nil
}

func (c *Cache) Zone(id string, worldID string) (types.Zone, error) {
	zone, err := c.localZone(id, worldID)
	if err != nil {
		return nil, err
	}

	if zone != nil {
		return zone, nil
	}

	return c.UpdateZone(id, worldID)
}

func (c *Cache) SetZone(zone types.Zone, worldID string) error {
	if err := c.db.SetZone(zone, worldID); err != nil {
		return err
	}

	if err := c.localSetZone(zone, worldID); err != nil {
		return err
	}

	return nil
}

func (c *Cache) UpdateZone(id string, worldID string) (types.Zone, error) {
	zone, err := c.db.GetZone(id, worldID)
	if err != nil {
		return nil, err
	}

	if zone != nil {
		err = c.localSetZone(zone, worldID)
	}

	return zone, err
}

func (c *Cache) localZone(id string, worldID string) (types.Zone, error) {
	c.zoneMutex.RLock()
	defer c.zoneMutex.RUnlock()

	zones, found := c.zones[worldID]
	if !found {
		return nil, nil
	}

	return zones[id], nil
}

func (c *Cache) localSetZone(zone types.Zone, worldID string) error {
	c.zoneMutex.Lock()
	defer c.zoneMutex.Unlock()

	_, found := c.zones[worldID]
	if !found {
		c.zones[worldID] = make(map[string]types.Zone)
	}

	c.zones[worldID][zone.ID()] = zone
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
	defer c.worldMutex.RUnlock()
	return c.worlds[id], nil
}

func (c *Cache) localSetWorld(world types.World) error {
	c.worldMutex.Lock()
	defer c.worldMutex.Unlock()
	c.worlds[world.ID()] = world
	return nil
}
