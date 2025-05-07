package subpub

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "poopsoup-app/pkg/pb/sub_pub"
)

func (i *Implementation) Publish(ctx context.Context, in *desc.PublishRequest) (*emptypb.Empty, error) {
	err := i.PubSub.Publish(in.GetKey(), in.GetData())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
