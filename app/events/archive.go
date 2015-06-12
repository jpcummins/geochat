package events

import (
	"github.com/jpcummins/geochat/app/types"
)

// Archive contains the latest events in a zone
type Archive struct {
	types.BaseEventData
	Events []*types.Event `json:"events"`
}

// Type implements EventData, which provides a hint to Event.Unmarshal on how
// to parse Event.Data
func (a *Archive) Type() string {
	return "archive"
}

// func newArchive(eventsJSON []string) (a *Archive) {
// 	a = &Archive{}
// 	events := list.New()
//
// 	// TODO: there must be a better way...
//
// 	for i := range eventsJSON {
// 		event := Event{}
// 		if err := json.Unmarshal([]byte(eventsJSON[i]), &event); err != nil {
// 			println(err.Error() + ": `" + eventsJSON[i] + "`")
// 			continue
// 		}
// 		println(string(i))
// 		events.PushBack(&event)
// 	}
//
// 	a.Events = make([]*Event, events.Len())
// 	i := 0
// 	for event := events.Front(); event != nil; event = event.Next() {
// 		a.Events[i] = event.Value.(*Event)
// 		i = i + 1
// 	}
// 	return
// }
