package chat

type Subscription struct {
	User   *User       `json:"user"`
	Events chan *Event `json:"-"`
	Zone   *Zone       `json:"-"`
}

func (s *Subscription) Unsubscribe() {
	for i, subscriber := range s.Zone.Subscribers {
		if subscriber == s {
			s.Zone.Subscribers = append(s.Zone.Subscribers[:i], s.Zone.Subscribers[i+1:]...)
			s.Zone.Publish(NewEvent(&Leave{User: s.User}))
			break
		}
	}
}
