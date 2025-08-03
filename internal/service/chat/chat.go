package chat

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type (
	Chat interface {
		// NewSubscriber creates a new subscriber and subscribes it to the chat.
		NewSubscriber(close func()) Subscriber

		// Subscribe adds a new subscriber to the chat.
		Subscribe(subscriber Subscriber)

		// Unsubscribe removes a subscriber from the chat.
		Unsubscribe(subscriber Subscriber)

		// Publish sends a message to all subscribers.
		Publish(ctx context.Context, message []byte) error
	}

	chat struct {
		// subscriberMessageBuffer is the size of the buffered channel for each subscriber.
		subscriberMessageBuffer int

		// publishLimiter is a rate limiter that limits the number of messages
		publishLimiter *rate.Limiter

		// logf is a function that logs messages.
		logf func(f string, v ...any)

		// subscribersMu is a mutex to protect the subscribers map.
		subscribersMu sync.Mutex
		subscribers   map[Subscriber]struct{}
	}
)

func NewChat(interval time.Duration, burst int) Chat {
	return &chat{
		subscriberMessageBuffer: 16,
		publishLimiter:          rate.NewLimiter(rate.Every(interval), burst),
		logf:                    func(f string, v ...any) {}, // Replace with actual logging function
		subscribers:             make(map[Subscriber]struct{}),
	}
}

// NewSubscriber creates a new subscriber and subscribes it to the chat.
func (c *chat) NewSubscriber(close func()) Subscriber {
	sub := NewSubscriber(close, c.subscriberMessageBuffer)
	c.Subscribe(sub)
	return sub
}

func (c *chat) Subscribe(subscriber Subscriber) {
	c.subscribersMu.Lock()
	defer c.subscribersMu.Unlock()

	c.subscribers[subscriber] = struct{}{}
}

func (c *chat) Unsubscribe(subscriber Subscriber) {
	c.subscribersMu.Lock()
	defer c.subscribersMu.Unlock()

	if _, exists := c.subscribers[subscriber]; exists {
		delete(c.subscribers, subscriber)
	}
}

func (c *chat) Publish(ctx context.Context, message []byte) error {
	c.subscribersMu.Lock()
	subscribers := make([]Subscriber, 0, len(c.subscribers))
	for sub := range c.subscribers {
		subscribers = append(subscribers, sub)
	}
	c.subscribersMu.Unlock()

	if err := c.publishLimiter.Wait(ctx); err != nil {
		return err
	}

	var closeSubs []Subscriber
	for _, sub := range subscribers {
		if err := sub.Send(message); err != nil {
			c.logf("failed to send message to subscriber: %v", err)
			closeSubs = append(closeSubs, sub)
			continue
		}
	}

	for _, sub := range closeSubs {
		c.Unsubscribe(sub)
		if err := sub.Close(); err != nil {
			c.logf("failed to close subscriber: %v", err)
		}
		c.logf("closed subscriber due to send error")
	}

	return nil
}
