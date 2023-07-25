package game

import (
	"fmt"
	"math/rand"
	"time"
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
	return hasPlayer(b.Players, player)
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
