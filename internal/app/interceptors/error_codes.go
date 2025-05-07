package interceptors

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"poopsoup-app/internal/domain"
)

func ErrorCodesUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrAlreadyExists) {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrInvalidArgument) {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrUnauthorized) {
		return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrForbidden) {
		return nil, status.Errorf(codes.PermissionDenied, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrTooManyRequests) {
		return nil, status.Errorf(codes.ResourceExhausted, "%s", err.Error())
	}

	log.Printf("[interceptor.Error] method: %s; error: %s", info.FullMethod, err.Error())
	return nil, status.Error(codes.Internal, "internal server error")
}

func ErrorCodesStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrNotFound) {
		return status.Errorf(codes.NotFound, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrAlreadyExists) {
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrInvalidArgument) {
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrUnauthorized) {
		return status.Errorf(codes.Unauthenticated, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrForbidden) {
		return status.Errorf(codes.PermissionDenied, "%s", err.Error())
	}
	if errors.Is(err, domain.ErrTooManyRequests) {
		return status.Errorf(codes.ResourceExhausted, "%s", err.Error())
	}

	log.Printf("[interceptor.Error] method: %s; error: %s", info.FullMethod, err.Error())
	return status.Error(codes.Internal, "internal server error")
}
