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

type BidCloseCode int

const (
	ContinueToGame BidCloseCode = iota
)

type BidSubscription struct {
	b     chan game.BidState
	close chan BidCloseCode
}

type BidManager struct {
	bidStates map[string]game.BidState
	subs      map[string](map[string]BidSubscription)
}

func (m BidManager) Get(id string) (game.BidState, bool) {
	b, exists := m.bidStates[id]
	return b, exists
}

func (m BidManager) Put(b game.BidState) {
	m.bidStates[b.GameID] = b
	subs, exists := m.subs[b.GameID]
	if !exists {
		m.subs[b.GameID] = make(map[string]BidSubscription)
		return
	}
	for _, s := range subs {
		s.b <- b
	}
}

func (m BidManager) Delete(id string, code BidCloseCode) {
	delete(m.bidStates, id)
	subs, exists := m.subs[id]
	if !exists {
		return
	}
	for _, s := range subs {
		s.close <- code
	}
	delete(m.subs, id)
}

func (m BidManager) Subscribe(id string, subscriber string) (BidSubscription, error) {
	_, exists := m.subs[id]
	if !exists {
		return BidSubscription{}, fmt.Errorf("Attempted to subscribe to a bid entry that doesn't exist")
	}
	sub := BidSubscription{make(chan game.BidState), make(chan BidCloseCode)}
	m.subs[id][subscriber] = sub
	return sub, nil
}

var bidManager = BidManager{
	bidStates: make(map[string]game.BidState),
	subs:      make(map[string]map[string]BidSubscription),
}

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

	lobbyManager.Delete(gameID, ContinueToBidding)

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
		bidManager.Delete(gameID, ContinueToGame)
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

	sub, err := bidManager.Subscribe(gameID, authInfo.Name)
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
			case b := <-sub.b:
				data, err := json.Marshal(b)
				if err != nil {
					log.Println("Got error when processing bid notification:", err)
					break
				}
				msg := fmt.Sprintf("event: update\ndata: %s\n\n", data)
				log.Printf("Sending message:\n%v", msg)
				fmt.Fprintf(w, msg)
			case code := <-sub.close:
				if code == ContinueToGame {
					log.Printf("Notifying of continue signal")
					fmt.Fprint(w, "event: continue\n\n")
				} else {
					log.Printf("Notifying of bidState deletion; code = %v", code)
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

func setupBidding(r fiber.Router) {
	r.Put("/:gameid", startBidding)
	r.Get("/:gameid", getBidState)
	r.Post("/:gameid/bid", submitBid)
	r.Get("/:gameid/subscribe", subscribeToBids)
}
