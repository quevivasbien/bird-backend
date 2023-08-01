package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
)

var lobbyManager = MakeManager[game.Lobby]()

func createLobby(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyID := c.Params("lobby")
	if _, exists := lobbyManager.Get(lobbyID); exists {
		return c.SendStatus(fiber.StatusConflict)
	}
	lobby := game.MakeLobby(lobbyID, authInfo.Name)

	lobbyManager.Put(lobby)
	return c.JSON(lobby)
}

func getLobbyState(c *fiber.Ctx) error {
	lobbyID := c.Params("lobby")
	lobby, exists := lobbyManager.Get(lobbyID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(lobby)
}

func subscribeToLobby(c *fiber.Ctx) error {
	lobbyID := c.Params("lobby")
	lobby, exists := lobbyManager.Get(lobbyID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of lobby in order to subscribe
	if !lobby.HasPlayer(authInfo.Name) && !authInfo.Admin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	err = lobbyManager.Subscribe(lobbyID, authInfo.Name, c)
	if err != nil {
		log.Println("When subscribing to lobby stream:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func swapLobbyOrder(c *fiber.Ctx) error {
	swap := struct {
		I int `json:"i"`
		J int `json:"j"`
	}{}
	if err := c.BodyParser(&swap); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyID := c.Params("lobby")
	lobby, exists := lobbyManager.Get(lobbyID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if !(lobby.Host == authInfo.Name || authInfo.Admin) {
		log.Printf("Attempted to swap lobby order with name %s, lobby host %s, and admin status = %v", authInfo.Name, lobby.Host, authInfo.Admin)
		return c.SendStatus(fiber.StatusForbidden)
	}
	i, j := swap.I, swap.J
	lobby.Players[i], lobby.Players[j] = lobby.Players[j], lobby.Players[i]
	lobbyManager.Put(lobby)
	return c.SendStatus(fiber.StatusAccepted)
}

func joinLobby(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyID := c.Params("lobby")
	lobby, exists := lobbyManager.Get(lobbyID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	for i, player := range lobby.Players {
		if player == "" {
			lobby.Players[i] = authInfo.Name
			lobbyManager.Put(lobby)
			return c.JSON(lobby)
		}
	}
	// lobby is full
	return c.SendStatus(fiber.StatusConflict)
}

func leaveLobby(c *fiber.Ctx) error {
	userInfo, err := UnloadTokenCookie(c)
	if err != nil || userInfo.Name == "" {
		c.SendStatus(fiber.StatusUnauthorized)
	}

	lobbyID := c.Params("lobby")
	lobby, exists := lobbyManager.Get(lobbyID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}

	for i, player := range lobby.Players {
		if player == userInfo.Name {
			lobby.Players[i] = ""
		}
	}

	// if player is host, set new host; delete game if no host remains
	lobby.Host = ""
	for _, player := range lobby.Players {
		if player != "" {
			lobby.Host = player
		}
	}
	if lobby.Host == "" {
		lobbyManager.Delete(lobbyID, EmptyCode)
		return c.SendStatus(fiber.StatusOK)
	} else {
		lobbyManager.Unsubscribe(lobbyID, userInfo.Name)
	}

	lobbyManager.Put(lobby)
	if err != nil {
		log.Println("When updating lobby players on db:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func setupLobbies(r fiber.Router) {
	r.Put("/:lobby", createLobby)
	r.Get("/:lobby", getLobbyState)
	r.Get("/:lobby/subscribe", subscribeToLobby)
	r.Post("/:lobby/swap", swapLobbyOrder)
	r.Post("/:lobby/join", joinLobby)
	r.Post("/:lobby/leave", leaveLobby)
}
