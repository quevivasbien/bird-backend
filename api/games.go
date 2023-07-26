package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
)

var gameManager = MakeManager[game.GameState]()

func setTrump(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	body := struct {
		Trump game.Color `json:"trump"`
	}{}
	if err = c.BodyParser(&body); err != nil {
		log.Println("Error parsing body of set trump request:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	game.Trump = body.Trump
	gameManager.Put(game)
	return c.SendStatus(fiber.StatusOK)
}

func getGameState(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(game)
}

func subscribeToGame(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	gameState, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of game in order to subscribe
	if !gameState.HasPlayer(authInfo.Name) && !authInfo.Admin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	_, err = gameManager.Subscribe(gameID, authInfo.Name, c)
	if err != nil {
		log.Println("When subscribing to bid stream:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func setupGames(r fiber.Router) {
	r.Get("/:gameid", getGameState)
	r.Post("/:gameid/trump", setTrump)
	r.Get("/:gameid/subscribe", subscribeToGame)
}
