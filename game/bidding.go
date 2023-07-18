package game

import "fmt"

type BidState struct {
	Done          bool    `json:"done"`
	Passed        [4]bool `json:"passed"`
	CurrentBidder int     `json:"currentBidder"`
	Bid           int     `json:"bid"`
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
	i := (b.CurrentBidder + 1) % 4
	for i != b.CurrentBidder {
		if !b.Passed[i] {
			b.CurrentBidder = i
			return
		}
	}
	b.Done = true
}

func (g GameState) ProcessBid(player string, amt int) error {
	// check that this is a valid bid
	if g.BidState.Done {
		return fmt.Errorf("Tried to send a bid while the game is not in the bidding stage")
	}
	playerIndex := -1
	for i, p := range g.Players {
		if player == p {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		return fmt.Errorf("Attempted to send a bid for a player who is not in this game")
	}
	if playerIndex != g.CurrentBidder {
		return fmt.Errorf("It is not currently this player's turn to bid")
	}
	if g.Passed[playerIndex] {
		return fmt.Errorf("Bidder already passed")
	}
	if amt < 0 || amt > 200 {
		return fmt.Errorf("Invalid bid amount")
	}

	if amt <= g.Bid {
		g.Passed[playerIndex] = true
	} else {
		g.Bid = amt
	}
	g.AdvanceBidder()

	return nil
}
