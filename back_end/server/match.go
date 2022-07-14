package server

import (
	"encoding/json"
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
	m.sendMatchBegin()
	firstTeamRound := "2"
	secondTeamRound := "2"
	round := "2"
	isFristRound := true
	banker := unknown
	for {
		game := &Game{
			Round:              round,
			FirstTeamRound:     firstTeamRound,
			SecondTeamRound:    secondTeamRound,
			IsFristRound:       isFristRound,
			Banker:             banker,
			PlayerConns:        m.PlayerConns,
			ShowMasterDoneChan: make(chan bool),
			ShowMasterChan:     make(chan *GameShowMasterRequest),
			BottomCardsChan:    make(chan []*Poker),
			PlayCardsChan: []chan []*Poker{
				make(chan []*Poker),
				make(chan []*Poker),
				make(chan []*Poker),
				make(chan []*Poker),
			},
		}
		m.sendGameToClients(game)
		gameResult := game.Run()
		if m.isFinish(game, gameResult) {
			log.Printf("比赛结束!, 获胜方: %d\n", banker)
			break
		}
		log.Println(gameResult)
		isFristRound = false
	}
}

func (m *Match) sendMatchBegin() {
	for _, conn := range m.PlayerConns {
		roomListMessage := ResponseMessage{
			MessageType: matchBegin,
		}
		sendMessage, _ := json.Marshal(roomListMessage)
		conn.WriteMessage(websocket.TextMessage, sendMessage)
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
