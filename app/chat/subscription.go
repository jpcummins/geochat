package chat

type Subscription struct {
	Id     string      `json:"id"`
	User   *User       `json:"user"`
	Events chan *Event `json:"-"`
	Zone   *Zone       `json:"-"`
}

func GetSubscription(id string) *Subscription {
	return subscriptions[id]
}
