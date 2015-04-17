package chat

import (
	"container/list"
	"encoding/json"
)

type Archive struct {
	Events []*Event `json:"events"`
}

func (a *Archive) Type() string {
	return "archive"
}

func newArchive(eventsJson []string) (a *Archive) {
	a = &Archive{}
	events := list.New()

	// TODO: there must be a better way...

	for i := range eventsJson {
		event := Event{}
		if err := json.Unmarshal([]byte(eventsJson[i]), &event); err != nil {
			println(err.Error() + ": `" + eventsJson[i] + "`")
			continue
		}
		events.PushBack(&event)
	}

	a.Events = make([]*Event, events.Len())
	i := 0
	for event := events.Front(); event != nil; event = event.Next() {
		a.Events[i] = event.Value.(*Event)
		i = i + 1
	}
	return
}
