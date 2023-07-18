package game

import "fmt"

type BidState struct {
	GameID        string    `json:"gameID"`
	Done          bool      `json:"done"`
	Players       [4]string `json:"players"`
	Hands         [4][]Card `json:"hands"`
	Widow         [5]Card   `json:"widow"`
	Passed        [4]bool   `json:"passed"`
	CurrentBidder int       `json:"currentBidder"`
	Bid           int       `json:"bid"`
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

func (b BidState) ProcessBid(player string, amt int) error {
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
