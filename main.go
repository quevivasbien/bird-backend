package main

import (
	"fmt"

	"github.com/quevivasbien/bird-backend/api"
	"github.com/quevivasbien/bird-backend/db"
)

const AWS_REGION = "us-east-1"

func main() {
	tables, _ := db.GetTables(AWS_REGION)
	tables.PutUser(db.User{
		Name:     "admin",
		Password: "admin",
		Admin:    true,
	})

	app, err := api.InitApp(AWS_REGION)
	if err != nil {
		panic(fmt.Sprintf("Error initializing app: %v", err))
	}
	app.Listen(":8081")
}
