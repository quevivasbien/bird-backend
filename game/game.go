package game

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
	Bid           int       `json:"bid"`
	BidWinner     int       `json:"bidWinner"`
}

// state of the game visible to a player during the game
type VisibleGameState struct {
	Hand          []Card `json:"hand"`
	Table         []Card `json:"table"`
	CurrentPlayer int    `json:"currentPlayer"`
	Trump         Color  `json:"trump"`
	Bid           int    `json:"bid"`
	BidWinner     int    `json:"bidWinner"`
}

func (g GameState) Visible(player int) VisibleGameState {
	return VisibleGameState{
		Hand:          g.Hands[player],
		Table:         g.Table,
		CurrentPlayer: g.CurrentPlayer,
		Trump:         g.Trump,
		Bid:           g.Bid,
		BidWinner:     g.BidWinner,
	}
}
