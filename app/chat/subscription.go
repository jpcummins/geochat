package chat

type Subscription struct {
	Id     string      `json:"id"`
	User   *User       `json:"user"`
	Events chan *Event `json:"-"`
	Zone   *Zone       `json:"-"`
}

func GetSubscription(id string) (s *Subscription, ok bool) {
	s, ok = subscriptions[id]

	if !ok {
		return nil, false
	}

	req := make(chan *Subscription)
	s.Zone.subscribe <- req // add channel to queue
	req <- s                // when ready, pass the subscription
	return <-req, true      // wait for processing to finish
}
