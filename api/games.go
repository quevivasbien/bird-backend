package api

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
	"github.com/quevivasbien/bird-game/game"
)

const LOBBY_TIMEOUT = time.Second * 30

func startBidding(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	gameID := c.Params("game")
	lobby, err := tables.GetLobby(gameID)
	if err != nil {
		var status int
		if _, ok := err.(db.ItemNotFound); ok {
			status = fiber.StatusNotFound
		} else {
			status = fiber.StatusInternalServerError
		}
		return c.SendStatus(status)
	}
	lobbyFull := true
	playerInLobby := false
	for _, player := range lobby.Players {
		if player == "" {
			lobbyFull = false
		}
		if player == authInfo.Name {
			playerInLobby = true
		}
	}
	if !lobbyFull {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if !playerInLobby {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	lobby.Started = true
	err = tables.PutLobby(lobby)
	if err != nil {
		log.Println("When sending started lobby back to db:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	go func() {
		time.Sleep(LOBBY_TIMEOUT)
		err := tables.DeleteLobby(gameID)
		if err != nil {
			log.Println("When attempting to delete started lobby:", err)
		}
	}()

	bidState := game.InitializeGame(gameID, lobby.Players)
	err = tables.PutBidState(bidState)
	if err != nil {
		log.Println("When putting new bid state in db:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func setupGames(r fiber.Router) {
	r.Post("/:game/bidding/start", startBidding)
}
