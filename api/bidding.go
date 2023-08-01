package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
	"github.com/quevivasbien/bird-game/utils"
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
	if lobby.Host != authInfo.Name {
		return c.SendStatus(fiber.StatusForbidden)
	}

	lobbyManager.Delete(gameID, ContinueCode)

	bidState := game.InitializeBidState(gameID, lobby.Players)
	bidManager.Put(bidState)

	return c.SendStatus(fiber.StatusOK)
}

func getBidState(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	gameID := c.Params("gameid")
	bidState, exists := bidManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	userIndex := utils.IndexOf(bidState.Players[:], authInfo.Name)
	if userIndex == -1 {
		log.Println("Tried to get game state for a player not in the game")
		return c.SendStatus(fiber.StatusForbidden)
	}
	return c.JSON(bidState.Visible(userIndex))
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
	if !bidState.HasPlayer(authInfo.Name) && !authInfo.Admin {
		return c.SendStatus(fiber.StatusForbidden)
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
		err = endBidding(bidState)
		if err != nil {
			log.Println("When ending bidding:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func endBidding(bidState game.BidState) error {
	if _, exists := bidManager.Get(bidState.ID); !exists {
		return fmt.Errorf("Tried to initialize a game from a BidState not in the bid manager")
	}
	game, err := bidState.InitGame()
	if err != nil {
		return fmt.Errorf("Error when initializing game from BidState: %v", err)
	}
	gameManager.Put(game)
	bidManager.Delete(bidState.ID, ContinueCode)
	return nil
}

func subscribeToBids(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	bidState, exists := bidManager.Get(gameID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of game in order to subscribe
	if !bidState.HasPlayer(authInfo.Name) {
		return c.SendStatus(fiber.StatusForbidden)
	}

	err = bidManager.Subscribe(gameID, authInfo.Name, c)
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
