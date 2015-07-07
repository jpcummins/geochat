package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/types"
	log "gopkg.in/inconshreveable/log15.v2"
	"sync"
)

type Worlds struct {
	sync.RWMutex
	db     types.DB
	pubsub types.PubSub
	worlds map[string]types.World
	logger log.Logger
}

func newWorlds(db types.DB, ps types.PubSub, logger log.Logger) *Worlds {
	return &Worlds{
		db:     db,
		pubsub: ps,
		worlds: make(map[string]types.World),
		logger: logger,
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

	if json == nil {
		return nil, nil
	}

	world := w.FromCache(id)
	if world == nil {
		if world, err = newWorld(id, w.db, w.pubsub, 10, w.logger); err != nil {
			return nil, err
		}
	}

	if err := world.Update(json); err != nil {
		return nil, err
	}

	return world, nil
}

func (w *Worlds) Save(world types.World) error {
	json, ok := world.PubSubJSON().(*types.WorldPubSubJSON)
	if !ok {
		return errors.New("Unable to serialize WorldPubSubJSON")
	}

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
