package pubsub

import (
	"context"
)

// Topic represents/identifies a particular Topic
type Topic string

// Broker pub-sub broker definition
type Broker interface {
	// Publish the given message onto the given topic
	Publish(ctx context.Context, m interface{}, topic Topic) error

	// Subscribe subscribes to the given topic
	Subscribe(ctx context.Context, subscriber *Subscriber, topic Topic) error

	// Topics retrieves a list of topics registered in this broker
	Topics(ctx context.Context) ([]Topic, error)

	// Shutdown shutsdown all subscribers gracefully
	Shutdown(ctx context.Context)
}
