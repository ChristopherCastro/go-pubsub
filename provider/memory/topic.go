package memory

import (
	"context"
	"sync"

	"github.com/ChristopherCastro/go-pubsub"
)

// topic defines a in-memory topic to which subscriber may subscribe to
type topic struct {
	id          pubsub.Topic
	subscribers map[string]*pubsub.Subscriber
	sync.RWMutex
}

type publishResult struct {
	subscriber *pubsub.Subscriber
	err        error
}

// newTopic creates a new topic
func newTopic(id pubsub.Topic) *topic {
	return &topic{
		id:          id,
		subscribers: map[string]*pubsub.Subscriber{},
	}
}

// publish sends the given message to each subscriber of this topic
func (t *topic) publish(ctx context.Context, m interface{}) []*publishResult {
	t.RLock()
	defer t.RUnlock()

	var wg sync.WaitGroup
	var errs sync.Map

	for _, s := range t.subscribers {
		wg.Add(1)
		subscriber := s

		go func() {
			defer wg.Done()

			result := &publishResult{
				subscriber: subscriber,
				err:        nil,
			}

			if err := subscriber.Deliver(ctx, m); err != nil {
				result.err = err
			}

			errs.Store(result, struct{}{})
		}()
	}

	wg.Wait()

	out := make([]*publishResult, 0)
	errs.Range(func(k, v interface{}) bool {
		out = append(out, k.(*publishResult))

		return true
	})

	return out
}

// subscribe attaches to this topic the given subscriber, attaching multiple times the same subscriber has no effects.
func (t *topic) subscribe(s *pubsub.Subscriber) {
	t.Lock()
	defer t.Unlock()

	t.subscribers[s.ID()] = s
}
