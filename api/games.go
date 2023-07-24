package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
)

type GameCloseCode int

const (
	FinishGame BidCloseCode = iota
)

type GameSubscription struct {
	b     chan game.GameState
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

func (m GameManager) Put(b game.GameState) {
	m.gameStates[b.GameID] = b
	subs, exists := m.subs[b.GameID]
	if !exists {
		m.subs[b.GameID] = make(map[string]GameSubscription)
		return
	}
	for _, s := range subs {
		s.b <- b
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

func setupGames(r fiber.Router) {
}
