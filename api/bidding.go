package api

import (
	"fmt"

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
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("When starting bids, attempted to fetch a lobby that doesn't exist")
	}
	if lobby.Host != authInfo.Name {
		c.Context().SetStatusCode(fiber.StatusForbidden)
		return c.SendString("You must be the lobby host to start bidding")
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
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("Requested bid state was not found in bid manager")
	}
	userIndex := utils.IndexOf(bidState.Players[:], authInfo.Name)
	if userIndex == -1 {
		c.Context().SetStatusCode(fiber.StatusForbidden)
		return c.SendString("Tried to get game state for a player not in the game")
	}
	return c.JSON(bidState.Visible(userIndex))
}

func submitBid(c *fiber.Ctx) error {
	gameID := c.Params("gameid")
	bidState, exists := bidManager.Get(gameID)
	if !exists {
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("Bid state was not found in bid manager")
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if !bidState.HasPlayer(authInfo.Name) {
		c.Context().SetStatusCode(fiber.StatusForbidden)
		return c.SendString("User is not a player in current game")
	}

	bid := struct {
		Amount int `json:"amount"`
	}{}
	if err := c.BodyParser(&bid); err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		return c.SendString(fmt.Sprintf("When parsing bid submission, got error %v", err))
	}

	err = bidState.ProcessBid(authInfo.Name, bid.Amount)
	if err != nil {
		c.Context().SetStatusCode(fiber.StatusBadRequest)
		return c.SendString(fmt.Sprintf("When checking if bid is valid for current bid state, got error %v", err))
	}

	bidManager.Put(bidState)
	if bidState.Done {
		err = endBidding(bidState)
		if err != nil {
			c.Context().SetStatusCode(fiber.StatusInternalServerError)
			return c.SendString(fmt.Sprintf("When ending bidding, got error %v", err))
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
		c.Context().SetStatusCode(fiber.StatusNotFound)
		return c.SendString("Requested bid state not found in bid manager")
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of game in order to subscribe
	if !bidState.HasPlayer(authInfo.Name) {
		c.Context().SetStatusCode(fiber.StatusForbidden)
		return c.SendString("User is not a player in current game")
	}

	err = bidManager.Subscribe(gameID, authInfo.Name, c)
	if err != nil {
		c.Context().SetStatusCode(fiber.StatusInternalServerError)
		return c.SendString(fmt.Sprintf("When subscribing to bid stream, got error %v", err))
	}

	return nil
}

func setupBidding(r fiber.Router) {
	r.Put("/:gameid", startBidding)
	r.Get("/:gameid", getBidState)
	r.Post("/:gameid", submitBid)
	r.Get("/:gameid/subscribe", subscribeToBids)
}
