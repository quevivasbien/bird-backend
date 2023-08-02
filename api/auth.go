package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
)

func loginHandler(c *fiber.Ctx) error {
	type LoginInput struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	var loginInput LoginInput
	if err := c.BodyParser(&loginInput); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var user db.User
	if tables == nil {
		user = db.User{
			Name:     loginInput.Name,
			Password: loginInput.Password,
		}
	} else {
		ok, u, err := tables.UserTable.ValidateUser(loginInput.Name, loginInput.Password)
		if !ok || err != nil {
			log.Println("When validating login:", err)
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		user = u
	}
	// login is ok; send jwt token
	userInfo, err := SetTokenCookie(c, user)
	if err != nil {
		log.Println("When getting JWT at login", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(userInfo)
}

func logoutHandler(c *fiber.Ctx) error {
	ClearTokenCookie(c)
	return c.SendStatus(fiber.StatusOK)
}

func createUserHandler(c *fiber.Ctx) error {
	if tables == nil {
		return c.SendStatus(fiber.StatusServiceUnavailable)
	}
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

func authStatusHandler(c *fiber.Ctx) error {
	userInfo, err := UnloadTokenCookie(c)
	if err != nil {
		log.Println("When unloading JWT cookie:", err)
	}
	return c.JSON(userInfo)
}

func setupAuth(r fiber.Router) {
	r.Post("/login", loginHandler)
	r.Post("/logout", logoutHandler)
	r.Post("/register", createUserHandler)
	r.Get("/status", authStatusHandler)
}
