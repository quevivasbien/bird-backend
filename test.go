package main

import (
	"fmt"

	"github.com/quevivasbien/bird-backend/api"
)

const AWS_REGION = "us-east-1"

func main() {
	app, err := api.InitApp(AWS_REGION)
	if err != nil {
		panic(fmt.Sprintf("Error initializing app: %v", err))
	}
	app.Listen(":8081")
}
