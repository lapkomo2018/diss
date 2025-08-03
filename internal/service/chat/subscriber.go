package chat

import (
	"errors"
	"sync"
)

type (
	Subscriber interface {
		Send(message []byte) error
		Messages() <-chan []byte
		Close() error
	}

	subscriber struct {
		msgs chan []byte

		closeOnce sync.Once
		close     func()
	}
)

func NewSubscriber(close func(), bufferSize int) Subscriber {
	return &subscriber{
		close: close,
		msgs:  make(chan []byte, bufferSize),
	}
}

func (s *subscriber) Send(message []byte) error {
	select {
	case s.msgs <- message:
		return nil
	default:
		return errors.New("subscriber buffer is full")
	}
}

func (s *subscriber) Messages() <-chan []byte {
	return s.msgs
}

// Close closes the subscriber's message channel.
func (s *subscriber) Close() error {
	s.closeOnce.Do(func() {
		close(s.msgs)
		if s.close != nil {
			s.close()
		}
	})
	return nil
}
