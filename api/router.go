package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-backend/db"
	"github.com/quevivasbien/bird-backend/game"
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
		err = SetTokenCookie(c, user)
		if err != nil {
			log.Println("When getting JWT at login", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func getLogoutHandler(tables db.Tables) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ClearTokenCookie(c)
		return c.SendStatus(fiber.StatusOK)
	}
}

func getCreateUserHandler(tables db.Tables) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type CreateUserInput struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		var input CreateUserInput
		if err := c.BodyParser(&input); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		alreadyExists, err := tables.UserExists(input.Name)
		if err != nil {
			log.Println("When checking if user exists:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if alreadyExists {
			return c.SendStatus(fiber.StatusConflict)
		}
		err = tables.PutUser(db.User{
			Name:     input.Name,
			Password: input.Password,
			Admin:    false,
		})
		if err != nil {
			log.Println("When creating new user on db:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusAccepted)
	}
}

func getCreateGameHandler(tables db.Tables) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authInfo, err := UnloadTokenCookie(c)
		if err != nil || authInfo.Name == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		lobby := db.Lobby{
			ID:      game.GetFreeGameID(),
			Host:    authInfo.Name,
			Players: [4]string{authInfo.Name},
		}
		err = tables.PutLobby(lobby)
		if err != nil {
			log.Println("When putting new lobby in db:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func getSubscribeToLobbyHandler(tables db.Tables) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}

func InitApi(r fiber.Router, region string) error {
	// app := fiber.New()
	// app.Use(
	// 	cors.New(cors.Config{
	// 		AllowOriginsFunc: func(origin string) bool {
	// 			return true
	// 		},
	// 		AllowHeaders:     "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	// 		AllowCredentials: true,
	// 	}),
	// )

	tables, err := db.GetTables("us-east-1")
	if err != nil {
		return fmt.Errorf("Error intializing tables: %v", err)
	}
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bird backend")
	})
	loginHandler := getLoginHandler(tables)
	r.Post("/login", loginHandler)
	r.Post("/logout", getLogoutHandler(tables))
	r.Post("/register", getCreateUserHandler(tables))
	r.Post("/games/create", getCreateGameHandler(tables))
	r.Get("/games/lobbies/:lobby", getSubscribeToLobbyHandler(tables))

	r.Get("/login/testAuth", func(c *fiber.Ctx) error {
		authInfo, err := UnloadTokenCookie(c)
		if err != nil {
			return c.SendString(fmt.Sprintf("Got error when unloading cookie: %v", err))
		}
		return c.SendString(fmt.Sprintf("Authinfo %v", authInfo))
	})

	return nil
}
