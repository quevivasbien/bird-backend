package game

type Lobby struct {
	ID      string    `json:"id"`
	Host    string    `json:"host"`
	Players [4]string `json:"players"`
}

func MakeLobby(id string, host string) Lobby {
	return Lobby{
		ID:      id,
		Host:    host,
		Players: [4]string{host},
	}
}

func (l Lobby) HasPlayer(player string) bool {
	return hasPlayer(l.Players, player)
}
