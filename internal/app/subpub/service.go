package subpub

import (
	desc "poopsoup-app/pkg/pb/sub_pub"
)

type PubSubService interface {
	Subscribe(subject string, handler MessageHandler) (Subscription, error)
	Publish(subject string, msg string) error
}

type Subscription interface {
	Unsubscribe()
}

type MessageHandler func(msg string)

type Implementation struct {
	desc.UnimplementedPubSubServer
	PubSub PubSubService
}

func New(ps PubSubService) *Implementation {
	return &Implementation{
		PubSub: ps,
	}
}
