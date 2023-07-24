package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
	"github.com/valyala/fasthttp"
)

type LobbyCloseCode int

const (
	ContinueToBidding LobbyCloseCode = iota
	LobbyEmpty
)

type LobbySubscription struct {
	l     chan game.Lobby
	close chan LobbyCloseCode
}

type LobbyManager struct {
	lobbies map[string]game.Lobby
	subs    map[string](map[string]LobbySubscription)
}

func (m LobbyManager) Get(id string) (game.Lobby, bool) {
	l, exists := m.lobbies[id]
	return l, exists
}

func (m LobbyManager) Put(l game.Lobby) {
	m.lobbies[l.ID] = l
	subs, exists := m.subs[l.ID]
	if !exists {
		m.subs[l.ID] = make(map[string]LobbySubscription)
		return
	}
	for _, s := range subs {
		s.l <- l
	}
}

func (m LobbyManager) Delete(id string, code LobbyCloseCode) {
	delete(m.lobbies, id)
	subs, exists := m.subs[id]
	if !exists {
		return
	}
	for _, s := range subs {
		s.close <- code
	}
	delete(m.subs, id)
}

func (m LobbyManager) Subscribe(id string, subscriber string) (LobbySubscription, error) {
	_, exists := m.subs[id]
	if !exists {
		return LobbySubscription{}, fmt.Errorf("Attempted to subscribe to a lobby entry that doesn't exist")
	}
	sub := LobbySubscription{make(chan game.Lobby), make(chan LobbyCloseCode)}
	m.subs[id][subscriber] = sub
	return sub, nil
}

func (m LobbyManager) Unsubscribe(id string, subscriber string) {
	_, exists := m.subs[id]
	if !exists {
		return
	}
	delete(m.subs[id], subscriber)
}

var lobbyManager = LobbyManager{
	lobbies: make(map[string]game.Lobby),
	subs:    make(map[string]map[string]LobbySubscription),
}

func createLobby(c *fiber.Ctx) error {
	authInfo, err := UnloadTokenCookie(c)
	if err != nil || authInfo.Name == "" {
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
	if err != nil || authInfo.Name == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// require player to be member of lobby in order to subscribe
	if !lobby.HasPlayer(authInfo.Name) && !authInfo.Admin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	sub, err := lobbyManager.Subscribe(lobbyID, authInfo.Name)
	if err != nil {
		log.Println("When subscribing to lobby stream:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		log.Println("Subscribed to lobby stream")
		for {
			select {
			case l := <-sub.l:
				data, err := json.Marshal(l)
				if err != nil {
					log.Println("Got error when processing lobby notification:", err)
					break
				}
				msg := fmt.Sprintf("event: update\ndata: %s\n\n", data)
				log.Printf("Sending message:\n%v", msg)
				fmt.Fprintf(w, msg)
			case code := <-sub.close:
				if code == ContinueToBidding {
					log.Printf("Notifying of continue signal")
					fmt.Fprint(w, "event: continue\n\n")
				} else {
					log.Printf("Notifying of lobby deletion; code = %v", code)
					fmt.Fprint(w, "event: delete\n\n")
				}
				return
			}
			err := w.Flush()
			if err != nil {
				log.Printf("Error while flushing: %v. Closing stream.", err)
				return
			}
		}
	}))

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
	if err != nil || authInfo.Name == "" {
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
	if err != nil || authInfo.Name == "" {
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
		lobbyManager.Delete(lobbyID, 1)
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
