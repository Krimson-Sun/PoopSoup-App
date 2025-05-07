package interceptors

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
)

func RecoverUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("[recovery] panic in: %s: %v \n %s", info.FullMethod, r, debug.Stack())
			log.Println(msg)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func RecoverStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("[recovery] panic in: %s: %v \n %s", info.FullMethod, r, debug.Stack())
			log.Println(msg)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(srv, ss)
}
