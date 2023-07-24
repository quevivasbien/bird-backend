package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/db"
	"github.com/valyala/fasthttp"
)

type LobbySubscription struct {
	l     chan db.Lobby
	close chan int
}

type LobbyManager struct {
	lobbies map[string]db.Lobby
	subs    map[string](map[string]LobbySubscription)
}

func (m LobbyManager) Get(id string) (db.Lobby, bool) {
	l, exists := m.lobbies[id]
	return l, exists
}

func (m LobbyManager) Put(l db.Lobby) {
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

func (m LobbyManager) Delete(id string) {
	delete(m.lobbies, id)
	subs, exists := m.subs[id]
	if !exists {
		return
	}
	for _, s := range subs {
		s.close <- 0
	}
	delete(m.subs, id)
}

func (m LobbyManager) Subscribe(id string, subscriber string) (LobbySubscription, error) {
	_, exists := m.subs[id]
	if !exists {
		return LobbySubscription{}, fmt.Errorf("Attempted to subscribe to a lobby entry that doesn't exist")
	}
	sub := LobbySubscription{make(chan db.Lobby), make(chan int)}
	m.subs[id][subscriber] = sub
	return sub, nil
}

var lobbyManager = LobbyManager{
	lobbies: make(map[string]db.Lobby),
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
	lobby := db.Lobby{
		ID:      lobbyID,
		Host:    authInfo.Name,
		Players: [4]string{authInfo.Name},
	}
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
	playerInLobby := false
	for _, p := range lobby.Players {
		if p == authInfo.Name {
			playerInLobby = true
			break
		}
	}
	if !playerInLobby && !authInfo.Admin {
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
			case <-sub.close:
				log.Println("Sending end message")
				fmt.Fprint(w, "event: end\n\n")
			}
			err := w.Flush()
			if err != nil {
				log.Printf("Error while flushing: %v. Closing stream.", err)
				break
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
	lobby, exists := lobbyManager.Get(lobbyID)
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
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
		lobbyManager.Delete(lobbyID)
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
	r.Get("/:lobby/subscribe", subscribeToLobby)
	r.Post("/:lobby/swap", swapLobbyOrder)
	r.Post("/:lobby/join", joinLobby)
	r.Post("/:lobby/leave", leaveLobby)
}
