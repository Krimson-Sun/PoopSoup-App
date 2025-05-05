package app

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"poopsoup-app/internal/app/subpub"
	desc "poopsoup-app/pkg/pb/sub_pub"
	"syscall"
)

type Options struct {
	grpcPort         int
	enableReflection bool
}

var defaultOptions = &Options{
	grpcPort:         50051,
	enableReflection: true,
}

type OptionsFunc func(*Options)

func WithGrpcPort(port int) OptionsFunc {
	return func(o *Options) {
		o.grpcPort = port
	}
}

func WithEnableReflection(enableReflection bool) OptionsFunc {
	return func(o *Options) {
		o.enableReflection = enableReflection
	}
}

type App struct {
	options *Options
	PubSub  subpub.PubSubService
}

func New(
	ps subpub.PubSubService,
	options ...OptionsFunc,
) *App {
	opts := defaultOptions
	for _, o := range options {
		o(opts)
	}
	return &App{
		options: opts,
		PubSub:  ps,
	}
}

func (a *App) Run() error {
	grpcEndpoint := fmt.Sprintf(":%d", a.options.grpcPort)
	srv := grpc.NewServer()
	service := subpub.New(a.PubSub)

	desc.RegisterPubSubServer(srv, service)

	if a.options.enableReflection {
		reflection.Register(srv)
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		log.Println("shutting down server...")

		srv.Stop()
	}()

	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("grpc server listening on port %d \n", a.options.grpcPort)

	// Start the server
	if err := srv.Serve(lis); err != nil {
		return err
	}

	log.Println("grpc server stopped")

	return nil
}
