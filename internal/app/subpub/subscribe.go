package subpub

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	desc "poopsoup-app/pkg/pb/sub_pub"
)

func (i *Implementation) Subscribe(in *desc.SubscribeRequest, stream grpc.ServerStreamingServer[desc.Event]) error {
	sub, err := i.PubSub.Subscribe(in.GetKey(), func(msg string) {
		stream.Send(&desc.Event{Data: msg})
	})
	if err != nil {
		return status.Errorf(codes.Unavailable, "Subscribe failed: %v", err)
	}
	<-stream.Context().Done()
	sub.Unsubscribe()
	return nil
}
