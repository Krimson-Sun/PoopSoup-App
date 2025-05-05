package subpub

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "poopsoup-app/pkg/pb/sub_pub"
)

func (i *Implementation) Publish(ctx context.Context, in *desc.PublishRequest) (*emptypb.Empty, error) {
	err := i.PubSub.Publish(in.GetKey(), in.GetData())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
