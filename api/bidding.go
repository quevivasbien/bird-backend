package api

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
)

var bidManager = MakeManager[game.BidState]()

func startBidding(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	gameID := c.Params("gameid")
	lobby, exists := lobbyManager.Get(gameID)
	if !exists {
		log.Println("When starting bids, attempted to fetch a lobby that doesn't exist")
		return c.SendStatus(fiber.StatusNotFound)
	}
	// check that player belongs in lobby and game is ready to start
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

	lobbyManager.Delete(gameID, ContinueCode)

	bidState := game.InitializeBidState(gameID, lobby.Players)
	bidManager.Put(bidState)

	return c.SendStatus(fiber.StatusOK)
}

func getBidState(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	bidState, exists := bidManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(bidState)
}

func submitBid(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	bidState, exists := bidManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if (authInfo.Name == "" || !bidState.HasPlayer(authInfo.Name)) && !authInfo.Admin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	bid := struct {
		Amount int `json:"amount"`
	}{}
	if err := c.BodyParser(&bid); err != nil {
		log.Println("When parsing bid submission:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = bidState.ProcessBid(authInfo.Name, bid.Amount)
	if err != nil {
		log.Println("When checking if bid is valid for current bid state:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	bidManager.Put(bidState)
	if bidState.Done {
		bidManager.Delete(gameID, ContinueCode)
	}

	// for testing
	if gameID == "test" {
		go func() {
			for strings.HasPrefix(bidState.Players[bidState.CurrentBidder], "dummy") {
				time.Sleep(time.Second)
				bidState.ProcessBid(bidState.Players[bidState.CurrentBidder], 0)
				bidManager.Put(bidState)
			}
			if bidState.Done {
				bidManager.Delete(gameID, ContinueCode)
			}
		}()
	}

	return c.SendStatus(fiber.StatusOK)
}

func subscribeToBids(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	bidState, exists := bidManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of game in order to subscribe
	if !bidState.HasPlayer(authInfo.Name) && !authInfo.Admin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	_, err = bidManager.Subscribe(gameID, authInfo.Name, c)
	if err != nil {
		log.Println("When subscribing to bid stream:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func setupBidding(r fiber.Router) {
	r.Put("/:gameid", startBidding)
	r.Get("/:gameid", getBidState)
	r.Post("/:gameid", submitBid)
	r.Get("/:gameid/subscribe", subscribeToBids)
}
