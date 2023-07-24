package game

type Lobby struct {
	ID      string    `json:"id"`
	Host    string    `json:"host"`
	Players [4]string `json:"players"`
}
