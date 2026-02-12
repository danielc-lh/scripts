package gcp

import (
	"context"
	"fmt"
	"os"
	"sync"

	pubsub "cloud.google.com/go/pubsub/v2"
)

type PubSubTransport struct {
	projectID     string
	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription
	client        *pubsub.Client
	lock          sync.RWMutex
}

func NewPubSubTransport(projectID string, topicNames []string) (*PubSubTransport, error) {
	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to create new pubsub client: %w", err)
	}

	// client.Publisher().Publish()
	// client.Subscriber().Receive()
}

func (t *PubSubTransport) Send(from, to string, data interface{}) {
	t.lock.RLock()
	topic, ok := t.topics[to]
	t.lock.RUnlock()
	if !ok {
		fmt.Fprintf(os.Stderr, "Topic %s not found\n", to)
		return
	}
	ctx := context.Background()
	msg := &pubsub.Message{
		Attributes: map[string]string{
			"from": from,
			"to":   to,
		},
		Data: []byte(fmt.Sprintf("%v", data)),
	}
	result := topic.Publish(ctx, msg)
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to publish: %v\n", err)
		return
	}
	fmt.Printf("Published message with ID: %s\n", id)
}

func (t *PubSubTransport) Receive(from, at string) interface{} {
	t.lock.RLock()
	sub, ok := t.subscriptions[at]
	t.lock.RUnlock()
	if !ok {
		fmt.Fprintf(os.Stderr, "Subscription %s not found\n", at)
		return nil
	}
	ctx := context.Background()
	var received interface{}
	cctx, cancel := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		if msg.Attributes["from"] == from {
			received = string(msg.Data)
			msg.Ack()
			cancel()
		}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to receive: %v\n", err)
		return nil
	}
	return received
}

func (t *PubSubTransport) Locations() []string {
	locs := make([]string, 0, len(t.topics))
	for k := range t.topics {
		locs = append(locs, k)
	}
	return locs
}
