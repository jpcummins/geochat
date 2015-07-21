package pubsub

import (
	"github.com/jpcummins/geochat/app/types"
)

const worldType types.PubSubEventType = "world"

type world struct {
	World *types.WorldPubSubJSON `json:"world"`
}

func World(json *types.WorldPubSubJSON) (*world, error) {
	w := &world{
		World: json,
	}
	return w, nil
}

func (w *world) Type() types.PubSubEventType {
	return worldType
}

func (w *world) BeforePublish(e types.PubSubEvent) error {
	return nil
}

func (w *world) OnReceive(e types.PubSubEvent) error {
	e.World().Update(w.World)
	return nil
}
