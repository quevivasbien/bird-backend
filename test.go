package main

import (
	"fmt"

	"github.com/quevivasbien/bird-backend/api"
)

func main() {
	app, err := api.InitApp()
	if err != nil {
		panic(fmt.Sprintf("Error initializing app: %v", err))
	}
	app.Listen(":8081")
}
