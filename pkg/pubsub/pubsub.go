package pubsub

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type MessageHandler[T interface{}] func(msg T)

type messageHandler[T interface{}] struct {
	HandlerFunc MessageHandler[T]
	ID          uuid.UUID
}

type Subscription interface {
	Unsubscribe()
}

type subscriptionImpl struct {
	unsub func()
}

func (s *subscriptionImpl) Unsubscribe() {
	s.unsub()
}

type SubPub[T interface{}] interface {
	Subscribe(subject string, handler MessageHandler[T]) (Subscription, error)
	Publish(subject string, msg T) error
	Close(ctx context.Context) error
}
type subPub[T interface{}] struct {
	subscribers map[string][]messageHandler[T]
	mtx         sync.RWMutex
	wg          sync.WaitGroup
}

func NewSubPub[T interface{}]() SubPub[T] {
	return &subPub[T]{
		subscribers: make(map[string][]messageHandler[T]),
		mtx:         sync.RWMutex{},
		wg:          sync.WaitGroup{},
	}
}

func (s *subPub[T]) Subscribe(subject string, handler MessageHandler[T]) (Subscription, error) {
	mh := messageHandler[T]{
		HandlerFunc: handler,
		ID:          uuid.New(),
	}
	s.mtx.Lock()
	s.subscribers[subject] = append(s.subscribers[subject], mh)
	s.mtx.Unlock()
	return &subscriptionImpl{
		unsub: func() {
			s.mtx.Lock()
			defer s.mtx.Unlock()
			handlers := s.subscribers[subject]
			for i, h := range handlers {
				if h.ID == mh.ID {
					s.subscribers[subject] = append(handlers[:i], handlers[i+1:]...)
					break
				}
			}
		},
	}, nil
}

func (s *subPub[T]) Publish(subject string, msg T) error {
	s.mtx.RLock()
	handlers, ok := s.subscribers[subject]
	s.mtx.RUnlock()
	if !ok || len(handlers) == 0 {
		return fmt.Errorf("no handler found for subject '%s'", subject)
	}

	for _, handler := range handlers {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			handler.HandlerFunc(msg)
		}()
	}
	return nil
}

func (s *subPub[T]) Close(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
