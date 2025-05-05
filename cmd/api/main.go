package main

import (
	"poopsoup-app/internal/app"
	"poopsoup-app/pkg/pubsub"
)

func main() {
	ps := pubsub.NewSubPub[string]()
	adapter := pubsub.NewAdapter(ps)
	App := app.New(adapter, app.WithEnableReflection(true))
	err := App.Run()
	if err != nil {
		panic(err)
	}
}
