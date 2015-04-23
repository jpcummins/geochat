package chat

type Subscription struct {
	Id     string      `json:"id"`
	User   *User       `json:"user"`
	Events chan *Event `json:"-"`
	Zone   *Zone       `json:"-"`
}

func GetSubscription(id string) (s *Subscription, ok bool) {
	s, ok = subscriptions[id]
	return
}
