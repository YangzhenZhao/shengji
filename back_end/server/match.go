package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Match struct {
	FirstTeamRound   string
	SecondTeamRound  string
	PlayerConns      []*websocket.Conn
	ReceiveGameChans []chan *Game
}

func (m *Match) Run() {
	firstTeamRound := "2"
	secondTeamRound := "2"
	round := "2"
	isFristRound := true
	banker := unknown
	for {
		resultChan := make(chan *GameResult)
		game := &Game{
			Round:           round,
			FirstTeamRound:  firstTeamRound,
			SecondTeamRound: secondTeamRound,
			IsFristRound:    isFristRound,
			Banker:          banker,
			PlayerConns:     m.PlayerConns,
			GameResultChan:  resultChan,
		}
		m.sendGameToClients(game)
		game.Run()
		gameResult := <-resultChan
		if m.isFinish(game, gameResult) {
			log.Printf("比赛结束!, 获胜方: %d\n", banker)
			break
		}
		log.Println(gameResult)
		isFristRound = false
	}

}

func (m *Match) isFinish(game *Game, gameResult *GameResult) bool {
	return game.Round == "A" && gameResult.FinalCardWinTeam == game.bankerTeam() && gameResult.Score < 80
}

func (m *Match) sendGameToClients(game *Game) {
	for _, receiveChan := range m.ReceiveGameChans {
		receiveChan <- game
	}
}
