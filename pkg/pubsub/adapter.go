package pubsub

import (
	"poopsoup-app/internal/app/subpub"
)

type Adapter struct {
	ps SubPub[string]
}

func NewAdapter(ps SubPub[string]) *Adapter {
	return &Adapter{ps: ps}
}

func (a *Adapter) Subscribe(subject string, handler subpub.MessageHandler) (subpub.Subscription, error) {
	pubsubHandler := func(msg string) {
		handler(msg)
	}

	sub, err := a.ps.Subscribe(subject, pubsubHandler)
	if err != nil {
		return nil, err
	}

	return &subscriptionAdapter{sub: sub}, nil
}

func (a *Adapter) Publish(subject string, msg string) error {
	return a.ps.Publish(subject, msg)
}

type subscriptionAdapter struct {
	sub Subscription
}

func (s *subscriptionAdapter) Unsubscribe() {
	s.sub.Unsubscribe()
}
