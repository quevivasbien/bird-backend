package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
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

func setupBidding(r fiber.Router) {
	r.Put("/:gameid", startBidding)
}
