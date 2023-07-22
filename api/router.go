package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
	"github.com/quevivasbien/bird-game/game"
)

var tables *db.Tables

func createGameHandler(c *fiber.Ctx) error {
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
	return c.JSON(lobby)
}

func subscribeToLobbyHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func InitApi(r fiber.Router, t db.Tables) error {
	tables = &t
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bird backend")
	})

	setupAuth(r.Group("/auth"))

	r.Post("/games/create", createGameHandler)
	r.Get("/games/lobbies/:lobby", subscribeToLobbyHandler)

	r.Get("/login/testAuth", func(c *fiber.Ctx) error {
		authInfo, err := UnloadTokenCookie(c)
		if err != nil {
			return c.SendString(fmt.Sprintf("Got error when unloading cookie: %v", err))
		}
		return c.SendString(fmt.Sprintf("Authinfo %v", authInfo))
	})

	return nil
}
