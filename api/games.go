package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
	"github.com/quevivasbien/bird-game/utils"
)

var gameManager = MakeManager[game.GameState]()

func getGameState(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("Requested game not found in game manager")
	}
	userIndex := utils.IndexOf(game.Players[:], authInfo.Name)
	if userIndex == -1 {
		c.Context().SetStatusCode(fiber.StatusForbidden)
		return c.SendString("Tried to get game state for a player not in the game")
	}
	return c.JSON(game.Visible(userIndex))
}

func getWidow(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("Requested game not found in game manager")
	}

	if authInfo.Name != game.Players[game.BidWinner] && !authInfo.Admin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.JSON(game.Widow)
}

// set trump and exchange cards with widow
func startRound(c *fiber.Ctx) error {
	_, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	body := struct {
		Trump     game.Color  `json:"trump"`
		ToWidow   []game.Card `json:"toWidow"`
		FromWidow []game.Card `json:"fromWidow"`
	}{}
	if err = c.BodyParser(&body); err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		return c.SendString(fmt.Sprintf("Error parsing body of start game request: %v", err))
	}
	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("Requested game not found in game manager")
	}
	err = game.ExchangeWithWidow(body.ToWidow, body.FromWidow)
	if err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		return c.SendString(fmt.Sprintf("When exchanging cards with widow, got error %v", err))
	}
	game.Trump = body.Trump
	gameManager.Put(game)
	return c.SendStatus(fiber.StatusOK)
}

func playCard(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	gameID := c.Params("gameid")
	gameState, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if !gameState.HasPlayer(authInfo.Name) {
		return c.SendStatus(fiber.StatusForbidden)
	}
	card := game.Card{}
	if err = c.BodyParser(&card); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	playerIndex := utils.IndexOf(gameState.Players[:], authInfo.Name)
	if playerIndex == -1 {
		return c.SendStatus(fiber.StatusForbidden)
	}
	err = gameState.PlayCard(playerIndex, card)
	if err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		return c.SendString(fmt.Sprintf("When trying to play card, got error %v", err))
	}
	gameManager.Put(gameState)
	return c.SendStatus(fiber.StatusOK)
}

func finishPlay(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	winner, err := game.FinishPlay()
	if err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		c.SendString(fmt.Sprintf("When attempting to finish play, got error %v", err))
	}
	gameManager.Put(game)
	out := struct {
		Winner int `json:"winner"`
	}{winner}
	return c.JSON(out)
}

func getScore(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	game, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	score0, score1, err := game.Score()
	if err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		return c.SendString(fmt.Sprintf("When trying to get game score, got error %v", err))
	}
	out := struct {
		Score0 int `json:"score0"`
		Score1 int `json:"score1"`
	}{score0, score1}
	return c.JSON(out)
}

func subscribeToGame(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	gameState, exists := gameManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of game in order to subscribe
	if !gameState.HasPlayer(authInfo.Name) {
		return c.SendStatus(fiber.StatusForbidden)
	}

	err = gameManager.Subscribe(gameID, authInfo.Name, c)
	if err != nil {
		log.Println("When subscribing to game stream:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func setupGames(r fiber.Router) {
	r.Get("/:gameid", getGameState)
	r.Get("/:gameid/widow", getWidow)
	r.Post("/:gameid/start", startRound)
	r.Post("/:gameid/play", playCard)
	r.Post("/:gameid/finish", finishPlay)
	r.Get("/:gameid/score", getScore)
	r.Get("/:gameid/subscribe", subscribeToGame)
}
