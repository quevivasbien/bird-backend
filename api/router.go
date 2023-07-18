package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-backend/db"
)

func getLoginHandler(tables db.Tables) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type LoginInput struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		var loginInput LoginInput
		if err := c.BodyParser(&loginInput); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		ok, user, err := tables.UserTable.ValidateUser(loginInput.Name, loginInput.Password)
		if !ok || err != nil {
			log.Println("When validating login:", err)
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		// login is ok; send jwt token
		jwt, expireTime, err := getToken(user)
		if err != nil {
			log.Println("When getting JWT at login")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    jwt,
			Expires:  expireTime,
			HTTPOnly: true,
			Secure:   true,
		})
		return c.SendStatus(fiber.StatusAccepted)
	}
}

func InitApp() (*fiber.App, error) {
	app := fiber.New()
	tables, err := db.GetTables("us-east-1")
	if err != nil {
		return app, fmt.Errorf("Error intializing tables: %v", err)
	}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bird backend")
	})
	app.Get("/login", getLoginHandler(tables))
	return app, nil
}
