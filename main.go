package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/quevivasbien/bird-game/api"
	"github.com/quevivasbien/bird-game/db"
	"github.com/quevivasbien/bird-game/template"
)

const AWS_REGION = "us-east-1"
const PORT = ":3000"

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	tables, err := db.GetTables(AWS_REGION)
	if err != nil {
		panic(fmt.Sprintf("Error initializing database tables: %v", err))
	}
	err = api.InitApi(app.Group("/api"), tables)
	if err != nil {
		panic(fmt.Sprintf("Error initializing API router: %v", err))
	}

	app.All(
		"/*",
		filesystem.New(filesystem.Config{
			Root:         template.GetBuild(),
			NotFoundFile: "index.html",
			Index:        "index.html",
		}),
	)

	if err := app.Listen(PORT); err != nil {
		panic(err)
	}
}
