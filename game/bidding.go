package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/quevivasbien/bird-game/utils"
)

type BidState struct {
	ID            string    `json:"id"`
	Done          bool      `json:"done"`
	Players       [4]string `json:"players"`
	Hands         [4][]Card `json:"hands"`
	Widow         [5]Card   `json:"widow"`
	Passed        [4]bool   `json:"passed"`
	CurrentBidder int       `json:"currentBidder"`
	Bid           int       `json:"bid"`
}

func (b BidState) GetID() string {
	return b.ID
}

func (b BidState) Visible(player int) interface{} {
	return VisibleBidState{
		ID:            b.ID,
		Done:          b.Done,
		Players:       b.Players,
		Hand:          b.Hands[player],
		Passed:        b.Passed,
		CurrentBidder: b.CurrentBidder,
		Bid:           b.Bid,
	}
}

func (b BidState) GetPlayers() []string {
	return b.Players[:]
}

func deal() ([4][]Card, [5]Card) {
	// get all cards
	allCards := []Card{Bird}
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
		rem := (i - 5 + 4) % 4
		hands[rem] = append(hands[rem], card)
	}
	return hands, widow
}

func InitializeBidState(id string, players [4]string) BidState {
	hands, widow := deal()
	return BidState{
		ID:      id,
		Players: players,
		Hands:   hands,
		Widow:   widow,
	}
}

func (b BidState) HasPlayer(player string) bool {
	return utils.Contains(b.Players[:], player)
}

func (b BidState) Winner() int {
	if !b.Done {
		return -1
	}
	return b.CurrentBidder
}

// find the next bidder who hasn't passed, and set them as the CurrentBidder
// if only one player remains, set bidding state to Done
func (b *BidState) AdvanceBidder() {
	numNotPassed := 0
	minDist := 0
	nextBidder := b.CurrentBidder
	for i, passed := range b.Passed {
		if i == b.CurrentBidder {
			continue
		}
		if !passed {
			numNotPassed++
			dist := (i - b.CurrentBidder + 4) % 4
			if dist < minDist || minDist == 0 {
				minDist = dist
				nextBidder = i
			}
		}
	}
	if numNotPassed == 0 || (b.Passed[b.CurrentBidder] && numNotPassed == 1) {
		b.Done = true
	}
	b.CurrentBidder = nextBidder
	if b.Players[b.CurrentBidder] == "" {
		b.setAIBid()
	}
}

func (b *BidState) ProcessBid(player string, amt int) error {
	// check that this is a valid bid
	if b.Done {
		return fmt.Errorf("Tried to send a bid while the game is not in the bidding stage")
	}
	playerIndex := -1
	for i, p := range b.Players {
		if player == p {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		return fmt.Errorf("Attempted to send a bid for a player who is not in this game")
	}
	if playerIndex != b.CurrentBidder {
		return fmt.Errorf("It is not currently this player's turn to bid")
	}
	if b.Passed[playerIndex] {
		return fmt.Errorf("Bidder already passed")
	}
	if amt < 0 || amt > 200 {
		return fmt.Errorf("Invalid bid amount")
	}

	if amt <= b.Bid {
		b.Passed[playerIndex] = true
	} else {
		b.Bid = amt
	}
	b.AdvanceBidder()

	return nil
}

func (b BidState) InitGame() (GameState, error) {
	if !b.Done {
		return GameState{}, fmt.Errorf("Cannot init game before bidding is done")
	}
	return GameState{
		ID:            b.ID,
		Players:       b.Players,
		Hands:         b.Hands,
		Widow:         b.Widow,
		CurrentPlayer: b.CurrentBidder,
		Bid:           b.Bid,
		BidWinner:     b.CurrentBidder,
		Table:         []Card{},
	}, nil
}

type VisibleBidState struct {
	ID            string    `json:"id"`
	Done          bool      `json:"done"`
	Players       [4]string `json:"players"`
	Hand          []Card    `json:"hand"`
	Passed        [4]bool   `json:"passed"`
	CurrentBidder int       `json:"currentBidder"`
	Bid           int       `json:"bid"`
}

func (b *BidState) setAIBid() {
	value := handValue(b.Hands[b.CurrentBidder])
	fmt.Printf("Player %d has hand value %d\n", b.CurrentBidder, value)
	if b.Bid < value {
		rem := value % 5
		b.Bid = value + (5 - rem)
	} else {
		b.Passed[b.CurrentBidder] = true
	}
	b.AdvanceBidder()
}

func handValue(h []Card) int {
	total := 0
	colorCounts := make(map[Color]int)
	colorCounts[Red] = 0
	colorCounts[Yellow] = 0
	colorCounts[Green] = 0
	colorCounts[Black] = 0
	for _, card := range h {
		if card.Value == 1 {
			total += 15
		} else if card.Value == 0 {
			colorCounts[Red]++
			colorCounts[Yellow]++
			colorCounts[Green]++
			colorCounts[Black]++
		} else {
			total += card.Value
			colorCounts[card.Color]++
		}
	}
	maxColorCount := 0
	for _, count := range colorCounts {
		if count > maxColorCount {
			maxColorCount = count
		}
	}
	total += maxColorCount * 5
	return total
}
