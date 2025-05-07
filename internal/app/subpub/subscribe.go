package subpub

import (
	"google.golang.org/grpc"
	desc "poopsoup-app/pkg/pb/sub_pub"
)

func (i *Implementation) Subscribe(in *desc.SubscribeRequest, stream grpc.ServerStreamingServer[desc.Event]) error {
	sub, err := i.PubSub.Subscribe(in.GetKey(), func(msg string) {
		stream.Send(&desc.Event{Data: msg})
	})
	if err != nil {
		return err
	}
	<-stream.Context().Done()
	sub.Unsubscribe()
	return nil
}
