package chat

type Subscriptions struct {
	subscriptions map[string]*Subscription

	// subscriptions i/o channels
	get chan (chan interface{})
	add chan (chan *Subscription)
}

func newSubscriptions() *Subscriptions {
	s := &Subscriptions{}

	s.subscriptions = make(map[string]*Subscription, 0)

	// Set up subscription i/o channels
	s.get = make(chan (chan interface{}))
	s.add = make(chan (chan *Subscription))

	// Handles subscription data structure I/O. Reads and writes to the
	// subscrptions hash table should only be preformed through this goroutine to
	// prevent deadlock or any other thread-unsafe outcome.
	go s.manage()
	return s
}

func (s *Subscriptions) manage() {
	for {
		select {
		case ch := <-s.get:
			id := (<-ch).(string)
			ch <- s.subscriptions[id]
			close(ch)
		case ch := <-s.add:
			subscription := <-ch
			s.subscriptions[subscription.id] = subscription
			close(ch)
		}
	}
}

func (s *Subscriptions) Get(id string) (subscription *Subscription, found bool) {
	ch := make(chan interface{})
	s.get <- ch
	ch <- id
	subscription = (<-ch).(*Subscription)
	found = subscription != nil
	return
}

func (s *Subscriptions) Add(subscription *Subscription) {
	ch := make(chan *Subscription)
	s.add <- ch
	ch <- subscription
	<-ch
	return
}
