package game

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
	GameID        string      `json:"gameID"`
	Players       [4]string   `json:"players"`
	Hands         [4][13]Card `json:"hands"`
	Discarded     [2][]Card   `json:"discarded"`
	Widow         [5]Card     `json:"widow"`
	Table         [4]Card     `json:"table"`
	CurrentPlayer int         `json:"currentPlayer"`
	Bid           int         `json:"bid"`
	BidWinner     int         `json:"bidWinner"`
	Bidding       bool        `json:"bidding"`
	Trump         Color       `json:"trump"`
}

// state of the game visible to a player during the game
type VisibleGameState struct {
	Hand          [13]Card `json:"hand"`
	Table         [4]Card  `json:"table"`
	CurrentPlayer int      `json:"currentPlayer"`
	Bid           int      `json:"bid"`
	BidWinner     int      `json:"bidWinner"`
	Bidding       bool     `json:"bidding"`
	Trump         Color    `json:"trump"`
}

func (g GameState) Visible(player int) VisibleGameState {
	return VisibleGameState{
		Hand:          g.Hands[player],
		Table:         g.Table,
		CurrentPlayer: g.CurrentPlayer,
		Bid:           g.Bid,
		BidWinner:     g.BidWinner,
		Bidding:       g.Bidding,
		Trump:         g.Trump,
	}
}
