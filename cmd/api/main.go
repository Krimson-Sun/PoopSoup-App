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

	appOpts := []app.OptionsFunc{}

	if raw, ok := os.LookupEnv("GRPC_PORT"); ok {
		port, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("invalid port: %s, using default %d \n", raw, 50051)
		} else {
			appOpts = append(appOpts, app.WithGrpcPort(port))
		}
	}

	if raw, ok := os.LookupEnv("ENABLE_REFLECTION"); ok {
		enableReflection, err := strconv.ParseBool(raw)
		if err != nil {
			log.Printf("invalid reflection value: %s, using default %t \n", raw, true)
		} else {
			appOpts = append(appOpts, app.WithEnableReflection(enableReflection))
		}
	}

	pubSubOpts := []pubsub.OptionsFunc{}

	if raw, ok := os.LookupEnv("WORKER_COUNT"); ok {
		workers, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("invalid worker count: %s, using default \n", raw)
		} else {
			pubSubOpts = append(pubSubOpts, pubsub.WithWorkerCount(workers))
		}
	}

	if raw, ok := os.LookupEnv("QUEUE_SIZE"); ok {
		queue, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("invalid queue size: %s, using default \n", raw)
		} else {
			pubSubOpts = append(pubSubOpts, pubsub.WithQueueSize(queue))
		}
	}

	ps := pubsub.NewSubPub[string]()
	App := app.New(ps, appOpts...)
	err := App.Run()
	if err != nil {
		panic(err)
	}
}
