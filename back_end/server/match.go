package server

import "github.com/gorilla/websocket"

type Match struct {
	FirstTeamRound   string
	SecondTeamRound  string
	PlayerConns      []*websocket.Conn
	ReceiveGameChans []chan *Game
}

func (m *Match) Run() {
	for {
		game := &Game{
			Round: "2",
		}
		m.sendGameToClients(game)
	}
}

func (m *Match) sendGameToClients(game *Game) {
	for _, receiveChan := range m.ReceiveGameChans {
		receiveChan <- game
	}
}
