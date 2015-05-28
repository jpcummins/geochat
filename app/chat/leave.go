package chat

type Leave struct {
	Subscriber *Subscription `json:"subscriber"`
}

func (m *Leave) Type() string {
	return "leave"
}

func (l *Leave) OnReceive(e *Event) error {
	// zone := l.Subscriber.GetZone()
	// subscribers := zone.GetSubscribers()
	//
	// for i, x := range subscribers {
	// 	if x.id == l.Subscriber.GetID() {
	// 		copy(subscribers[i:], subscribers[i+1:])
	// 		subscribers[len(subscribers)-1] = nil
	// 		subscribers = subscribers[:len(subscribers)-1]
	// 		decrementZoneSubscriptionCounts(zone) // bubble up the count
	// 		return nil
	// 	}
	// }
	return nil
}
