package chat

import (
	"github.com/jpcummins/geochat/app/types"
	"sync"
)

type Worlds struct {
	sync.RWMutex
	db     types.DB
	worlds map[string]types.World
}

func newWorlds(db types.DB) *Worlds {
	return &Worlds{
		db:     db,
		worlds: make(map[string]types.World),
	}
}

func (w *Worlds) World(id string) (types.World, error) {
	if cachedWorld := w.FromCache(id); cachedWorld != nil {
		return cachedWorld, nil
	}

	return w.FromDB(id)
}

func (w *Worlds) FromCache(id string) types.World {
	w.RLock()
	defer w.RUnlock()
	return w.worlds[id]
}

func (w *Worlds) FromDB(id string) (types.World, error) {
	json, err := w.db.World(id)
	if err != nil {
		return nil, err
	}

	world := w.FromCache(id)
	if world == nil {
		return nil, nil
	}

	world.Update(json)
	w.updateCache(world)
	return world, nil
}

func (w *Worlds) Save(world types.World) error {
	json := world.ServerJSON()
	if err := w.db.SaveWorld(json); err != nil {
		return err
	}

	w.updateCache(world)
	return nil
}

func (w *Worlds) updateCache(world types.World) {
	w.Lock()
	defer w.Unlock()
	w.worlds[world.ID()] = world
}
