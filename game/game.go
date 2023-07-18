package game

import (
	"fmt"
	"math/rand"
	"time"
)

const GAME_ID_LENGTH = 8

type Color int

const (
	Red Color = iota + 1
	Yellow
	Green
	Black
)

type Card struct {
	Color Color
	Value int
}

var Rook Card = Card{0, 0}

type GameState struct {
	GameID        string    `json:"gameID"`
	Players       [4]string `json:"players"`
	Hands         [4][]Card `json:"hands"`
	Discarded     [2][]Card `json:"discarded"`
	Widow         [5]Card   `json:"widow"`
	Table         []Card    `json:"table"`
	CurrentPlayer int       `json:"currentPlayer"`
	Trump         Color     `json:"trump"`
	BidState      `json:"bidState"`
}

// state of the game visible to a player during the game
type VisibleGameState struct {
	Hand          []Card `json:"hand"`
	Table         []Card `json:"table"`
	CurrentPlayer int    `json:"currentPlayer"`
	Trump         Color  `json:"trump"`
	BidState      `json:"bidState"`
}

func (g GameState) Visible(player int) VisibleGameState {
	return VisibleGameState{
		Hand:          g.Hands[player],
		Table:         g.Table,
		CurrentPlayer: g.CurrentPlayer,
		Trump:         g.Trump,
		BidState:      g.BidState,
	}
}

func getFreeGameID() string {
	// todo: check for duplicates?
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, GAME_ID_LENGTH+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b[2:GAME_ID_LENGTH+2])
}

func deal() ([4][]Card, [5]Card) {
	// get all cards
	allCards := []Card{Rook}
	for suite := Red; suite <= Black; suite++ {
		for value := 1; value <= 14; value++ {
			allCards = append(allCards, Card{suite, value})
		}
	}
	// get random indices and distribute cards
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(allCards))
	hands := [4][]Card{}
	widow := [5]Card{}
	for i, j := range perm {
		card := allCards[j]
		if i < 5 {
			widow[i] = card
			continue
		}
		// div := (i - 5) / 4
		rem := (i - 5) % 4
		hands[rem] = append(hands[rem], card)
	}
	return hands, widow
}

func InitializeGame(players [4]string) GameState {
	gameID := getFreeGameID()
	hands, widow := deal()
	return GameState{
		GameID:  gameID,
		Players: players,
		Hands:   hands,
		Widow:   widow,
	}
}
