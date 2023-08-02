package game

import "github.com/quevivasbien/bird-game/utils"

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

func (l Lobby) GetID() string {
	return l.ID
}

func (l Lobby) Visible(int) interface{} {
	return l
}

func (l Lobby) GetPlayers() []string {
	return l.Players[:]
}

func (l Lobby) HasPlayer(player string) bool {
	return utils.Contains(l.Players[:], player)
}
