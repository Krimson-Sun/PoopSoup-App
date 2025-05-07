package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"poopsoup-app/internal/app"
	"poopsoup-app/pkg/pubsub"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	opts := []app.OptionsFunc{}

	if raw, ok := os.LookupEnv("GRPC_PORT"); ok {
		port, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("invalid port: %s, using default %d \n", raw, 50051)
		} else {
			opts = append(opts, app.WithGrpcPort(port))
		}
	}

	if raw, ok := os.LookupEnv("ENABLE_REFLECTION"); ok {
		enableReflection, err := strconv.ParseBool(raw)
		if err != nil {
			log.Printf("invalid reflection value: %s, using default %t \n", raw, true)
		} else {
			opts = append(opts, app.WithEnableReflection(enableReflection))
		}
	}

	ps := pubsub.NewSubPub[string]()
	adapter := pubsub.NewAdapter(ps)
	App := app.New(adapter, opts...)
	err := App.Run()
	if err != nil {
		panic(err)
	}
}
