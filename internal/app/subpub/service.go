package subpub

import (
	desc "poopsoup-app/pkg/pb/sub_pub"
	"poopsoup-app/pkg/pubsub"
)

type Implementation struct {
	desc.UnimplementedPubSubServer
	PubSub pubsub.SubPub[string]
}

func New(ps pubsub.SubPub[string]) *Implementation {
	return &Implementation{
		PubSub: ps,
	}
}
