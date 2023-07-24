package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
)

func createLobby(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyID := c.Params("lobby")
	_, err = tables.GetLobby(lobbyID)
	if err != nil {
		if _, ok := err.(db.ItemNotFound); !ok {
			log.Println("When checking if lobbyID is taken in an active lobby:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else {
		return c.SendStatus(fiber.StatusConflict)
	}
	_, err = tables.GetGameState(lobbyID)
	if err != nil {
		if _, ok := err.(db.ItemNotFound); !ok {
			log.Println("When checking if lobbyID is taken in an active game:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else {
		return c.SendStatus(fiber.StatusConflict)
	}
	lobby := db.Lobby{
		ID:      lobbyID,
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
	lobbyMap, err := dbCache.Get(tables.LobbyTable, c.Params("lobby"))
	if err != nil {
		log.Println("When getting lobby state from cache", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if lobbyMap == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	lobby, err := db.LobbyFromItemMap(lobbyMap)
	if err != nil {
		log.Println("When unmarshalling cached lobby state", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(lobby)
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
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyMap, err := dbCache.Get(tables.LobbyTable, c.Params("lobby"))
	if err != nil {
		if _, ok := err.(db.ItemNotFound); ok {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Println("When getting lobby state from cache", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	lobby, err := db.LobbyFromItemMap(lobbyMap)
	if err != nil {
		log.Println("When unmarshalling cached lobby state", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !(lobby.Host == authInfo.Name || authInfo.Admin) {
		return c.SendStatus(fiber.StatusForbidden)
	}
	i, j := swap.I, swap.J
	lobby.Players[i], lobby.Players[j] = lobby.Players[j], lobby.Players[i]
	err = tables.PutLobby(lobby)
	if err != nil {
		log.Println("When putting updated lobby in db:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func joinLobby(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	lobbyID := c.Params("lobby")
	lobbyMap, err := dbCache.Get(tables.LobbyTable, lobbyID)
	if err != nil {
		log.Println("When getting lobby state from cache", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if lobbyMap == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	lobby, err := db.LobbyFromItemMap(lobbyMap)
	if err != nil {
		log.Println("When unmarshalling cached lobby state", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	for i, player := range lobby.Players {
		if player == "" {
			lobby.Players[i] = authInfo.Name
			tables.PutLobby(lobby)
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
	lobby, err := tables.GetLobby(lobbyID)
	if err != nil {
		log.Println("When fetching lobby state from db:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for i, player := range lobby.Players {
		if player == userInfo.Name {
			lobby.Players[i] = ""
		}
	}

	// if player is host, set new host, or delete game if no host remains
	lobby.Host = ""
	for _, player := range lobby.Players {
		if player != "" {
			lobby.Host = player
		}
	}
	if lobby.Host == "" {
		tables.DeleteLobby(lobby.ID)
		return c.SendStatus(fiber.StatusOK)
	}

	err = tables.PutLobby(lobby)
	if err != nil {
		log.Println("When updating lobby players on db:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func setupLobbies(r fiber.Router) {
	r.Put("/:lobby", createLobby)
	r.Get("/:lobby", getLobbyState)
	r.Post("/:lobby/swap", swapLobbyOrder)
	r.Post("/:lobby/join", joinLobby)
	r.Post("/:lobby/leave", leaveLobby)
}
