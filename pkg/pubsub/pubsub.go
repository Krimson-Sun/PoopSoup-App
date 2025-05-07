package pubsub

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"poopsoup-app/internal/domain"
	"sync"
	"time"
)

type Options struct {
	WorkerCount int
	QueueSize   int
}

var defaultOptions = &Options{
	WorkerCount: 5,
	QueueSize:   100,
}

type OptionsFunc func(*Options)

func WithWorkerCount(count int) OptionsFunc {
	return func(o *Options) {
		o.WorkerCount = count
	}
}

func WithQueueSize(size int) OptionsFunc {
	return func(o *Options) {
		o.QueueSize = size
	}
}

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
	jobs        chan job[T]
}

func NewSubPub[T interface{}](opts ...OptionsFunc) SubPub[T] {
	options := *defaultOptions

	for _, opt := range opts {
		opt(&options)
	}
	sp := &subPub[T]{
		subscribers: make(map[string][]messageHandler[T]),
		jobs:        make(chan job[T], options.QueueSize),
	}

	for i := 0; i < options.WorkerCount; i++ {
		go sp.worker()
	}
	return sp
}

func (s *subPub[T]) worker() {
	for j := range s.jobs {
		j.mh.HandlerFunc(j.msg)
		s.wg.Done()
	}
}

type job[T any] struct {
	mh  messageHandler[T]
	msg T
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
		return fmt.Errorf("no handler found for subject '%s': %w", subject, domain.ErrNotFound)
	}

	for _, handler := range handlers {
		s.wg.Add(1)
		s.jobs <- job[T]{
			mh:  handler,
			msg: msg,
		}
	}
	return nil
}

func (s *subPub[T]) Close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	close(s.jobs)
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
