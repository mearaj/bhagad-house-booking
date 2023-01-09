package service

import (
	"errors"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"sync"
)

type Topic int

const (
	TopicBookingsFetched Topic = iota
	TopicSearchBookings
	TopicUserLoggedInOut
	TopicCreateBooking
	TopicUpdateBooking
	TopicDeleteBooking
)

var AllTopicsArr = [...]Topic{
	TopicBookingsFetched,
	TopicSearchBookings,
	TopicUserLoggedInOut,
	TopicCreateBooking,
	TopicUpdateBooking,
	TopicDeleteBooking,
}

type BookingChangedEventData struct{}
type CustomersChangeEventData struct{ BookingPublicKey string }

type Event struct {
	Data  interface{}
	Topic Topic
}

type EventCallback func(event Event)

type subscriber struct {
	events        chan Event
	topics        utils.Map[Topic, struct{}]
	closed        bool
	closedMutex   sync.RWMutex
	callback      EventCallback
	callbackMutex sync.RWMutex
}
type Subscriber interface {
	// Events the channel where Event can be received,
	Events() <-chan Event
	// Subscribe should subscribe to all events if empty, can return error esp when subscription is closed
	Subscribe(...Topic) error
	// IsSubscribedTo returns bool indicating subscriber's subscription status to a topic
	//  can return error esp when subscription is closed
	IsSubscribedTo(topic Topic) (bool, error)
	// UnSubscribe should unsubscribe to all events if empty, the Subscriber should still be open/usable
	//  can return error esp when subscription is closed
	UnSubscribe(...Topic) error
	// Close should make this interface unusable and release resources
	Close()
	// IsClosed indicates this interface is still usable
	IsClosed() bool
	// Topics returns topics to which subscriber is subscribed to,
	// can return error esp when subscription is closed
	Topics() ([]Topic, error)
	// SubscribeWithCallback (optionally listening with callback),
	// can return error esp when subscription is closed
	SubscribeWithCallback(callback EventCallback)
}

func newSubscriber() *subscriber {
	return &subscriber{
		events: make(chan Event, 100),
		closed: false,
		topics: utils.NewMap[Topic, struct{}](),
	}
}

func (s *subscriber) Events() <-chan Event {
	s.closedMutex.RLock()
	defer s.closedMutex.RUnlock()
	return s.events
}

func (s *subscriber) isError() error {
	if s.IsClosed() {
		return errors.New("subscription already closed")
	}
	return nil
}

// Subscribe subscribes to all events if empty
func (s *subscriber) Subscribe(topics ...Topic) (err error) {
	if err = s.isError(); err != nil {
		return err
	}
	if len(topics) == 0 {
		for _, eachTopic := range AllTopicsArr {
			s.topics.Set(eachTopic, struct{}{})
		}
		return nil
	}

	for _, eachTopic := range topics {
		s.topics.Set(eachTopic, struct{}{})
	}
	return nil
}

func (s *subscriber) IsSubscribedTo(topic Topic) (ok bool, err error) {
	if err = s.isError(); err != nil {
		return ok, err
	}
	_, ok = s.topics.Get(topic)
	return ok, err
}

// UnSubscribe unsubscribes to all events if empty, subscriber is still referenced
func (s *subscriber) UnSubscribe(topics ...Topic) (err error) {
	if err = s.isError(); err != nil {
		return err
	}
	if len(topics) == 0 {
		s.topics.Clear()
		return err
	}

	for _, eachTopic := range topics {
		s.topics.Delete(eachTopic)
	}
	return err
}

// Close closes the Event chan and UnSubscribes to all events and clears EventCallback
func (s *subscriber) Close() {
	_ = s.UnSubscribe()
	s.closedMutex.Lock()
	s.closed = true
	s.closedMutex.Unlock()
	s.callbackMutex.Lock()
	s.callback = nil
	s.callbackMutex.Unlock()
	go close(s.events)
}
func (s *subscriber) IsClosed() bool {
	s.closedMutex.RLock()
	defer s.closedMutex.RUnlock()
	return s.closed
}

func (s *subscriber) Topics() (topics []Topic, err error) {
	if err = s.isError(); err != nil {
		return nil, err
	}
	return s.topics.Keys(), err
}
func (s *subscriber) fire(event Event) {
	if ok, err := s.IsSubscribedTo(event.Topic); ok && err == nil && !s.IsClosed() {
		select {
		case s.events <- event:
		default:
			// channel buffer is full, empty first element and append to last element
			select {
			case <-s.Events():
				s.events <- event
			default:
			}
		}
		s.callbackMutex.RLock()
		if s.callback != nil {
			s.callback(event)
		}
		s.callbackMutex.RUnlock()
	}
}
func (s *subscriber) SubscribeWithCallback(callback EventCallback) {
	s.callbackMutex.Lock()
	defer s.callbackMutex.Unlock()
	s.callback = callback
}

type subscribers = utils.Map[*subscriber, struct{}]

type eventBroker struct {
	cachedEvents utils.Map[Topic, Event]
	subscribers  subscribers
}

func newEventBroker() *eventBroker {
	return &eventBroker{
		cachedEvents: utils.NewMap[Topic, Event](),
		subscribers:  utils.NewMap[*subscriber, struct{}](),
	}
}

func (eb *eventBroker) addSubscriber(sub *subscriber) {
	eb.subscribers.Set(sub, struct{}{})
	for _, e := range eb.cachedEvents.Values() {
		if ok, err := sub.IsSubscribedTo(e.Topic); ok && err == nil {
			go sub.fire(e)
		}
	}
}

func (eb *eventBroker) Fire(event Event) {
	//eb.cachedEventsMutex.Lock()
	//defer eb.cachedEventsMutex.Unlock()
	for _, sub := range eb.subscribers.Keys() {
		eb.cachedEvents.Set(event.Topic, event)
		if sub.IsClosed() {
			eb.subscribers.Delete(sub)
			continue
		}
		go sub.fire(event)
	}
}
