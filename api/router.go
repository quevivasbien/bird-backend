package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
	"github.com/quevivasbien/bird-game/game"
)

const CACHE_UPDATE_INTVL time.Duration = time.Millisecond * 500
const CACHE_FLUSH_INTVL time.Duration = time.Second * 30

var tables *db.Tables
var dbCache = db.MakeCache(CACHE_UPDATE_INTVL, CACHE_FLUSH_INTVL)

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

func getLobbyState(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyMap, err := dbCache.Get(tables.LobbyTable, c.Params("lobby"))
	if err != nil {
		log.Println("When getting lobby state from cache", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	lobby, err := db.LobbyFromItemMap(lobbyMap)
	if err != nil {
		log.Println("When unmarshalling cached lobby state", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(lobby)
}

func InitApi(r fiber.Router, t db.Tables) error {
	tables = &t
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bird backend")
	})

	setupAuth(r.Group("/auth"))

	r.Post("/games/create", createGameHandler)
	r.Get("/lobbies/:lobby", getLobbyState)

	r.Get("/login/testAuth", func(c *fiber.Ctx) error {
		authInfo, err := UnloadTokenCookie(c)
		if err != nil {
			return c.SendString(fmt.Sprintf("Got error when unloading cookie: %v", err))
		}
		return c.SendString(fmt.Sprintf("Authinfo %v", authInfo))
	})

	return nil
}
