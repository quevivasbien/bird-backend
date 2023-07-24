package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
)

const CACHE_UPDATE_INTVL time.Duration = time.Millisecond * 500
const CACHE_FLUSH_INTVL time.Duration = time.Second * 30

var tables *db.Tables

func InitApi(r fiber.Router, t db.Tables) error {
	tables = &t
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bird backend")
	})

	setupAuth(r.Group("/auth"))
	setupLobbies(r.Group("/lobbies"))
	setupBidding(r.Group("/bidding"))
	setupGames(r.Group("/games"))

	r.Get("/login/testAuth", func(c *fiber.Ctx) error {
		authInfo, err := UnloadTokenCookie(c)
		if err != nil {
			return c.SendString(fmt.Sprintf("Got error when unloading cookie: %v", err))
		}
		return c.SendString(fmt.Sprintf("Authinfo %v", authInfo))
	})

	return nil
}
