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

type GameCloseCode int

const (
	FinishGame GameCloseCode = iota
)

type GameSubscription struct {
	g     chan game.GameState
	close chan GameCloseCode
}

type GameManager struct {
	gameStates map[string]game.GameState
	subs       map[string](map[string]GameSubscription)
}

func (m GameManager) Get(id string) (game.GameState, bool) {
	b, exists := m.gameStates[id]
	return b, exists
}

func (m GameManager) Put(g game.GameState) {
	m.gameStates[g.ID] = g
	subs, exists := m.subs[g.ID]
	if !exists {
		m.subs[g.ID] = make(map[string]GameSubscription)
		return
	}
	for _, s := range subs {
		s.g <- g
	}
}

func (m GameManager) Delete(id string, code GameCloseCode) {
	delete(m.gameStates, id)
	subs, exists := m.subs[id]
	if !exists {
		return
	}
	for _, s := range subs {
		s.close <- code
	}
	delete(m.subs, id)
}

func (m GameManager) Subscribe(id string, subscriber string) (GameSubscription, error) {
	_, exists := m.subs[id]
	if !exists {
		return GameSubscription{}, fmt.Errorf("Attempted to subscribe to a bid entry that doesn't exist")
	}
	sub := GameSubscription{make(chan game.GameState), make(chan GameCloseCode)}
	m.subs[id][subscriber] = sub
	return sub, nil
}

var gameManager = GameManager{
	gameStates: make(map[string]game.GameState),
	subs:       make(map[string]map[string]GameSubscription),
}

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

	sub, err := gameManager.Subscribe(gameID, authInfo.Name)
	if err != nil {
		log.Println("When subscribing to bid stream:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		log.Println("Subscribed to bid stream")
		for {
			select {
			case g := <-sub.g:
				data, err := json.Marshal(g)
				if err != nil {
					log.Println("Got error when processing bid notification:", err)
					break
				}
				msg := fmt.Sprintf("event: update\ndata: %s\n\n", data)
				log.Printf("Sending message:\n%v", msg)
				fmt.Fprintf(w, msg)
			case code := <-sub.close:
				if code == FinishGame {
					log.Printf("Notifying of continue signal")
					fmt.Fprintf(w, "event: continue\ndata: %d\n\n", code)
				} else {
					log.Printf("Notifying of gameState deletion; code = %v", code)
					fmt.Fprintf(w, "event: delete\ndata: %d\n\n", code)
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

func setupGames(r fiber.Router) {
	r.Get("/:gameid", getGameState)
	r.Post("/:gameid/trump", setTrump)
	r.Get("/:gameid/subscribe", subscribeToGame)
}
