package server

import (
	"encoding/json"
	"log"

	"github.com/YangzhenZhao/shengji/back_end/server/common"
	"github.com/YangzhenZhao/shengji/back_end/server/dto"
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
	gameBanker := unknown
	for {
		game := &Game{
			Round:              round,
			FirstTeamRound:     firstTeamRound,
			SecondTeamRound:    secondTeamRound,
			IsFristRound:       isFristRound,
			Banker:             gameBanker,
			PlayerConns:        m.PlayerConns,
			ShowMasterDoneChan: make(chan bool),
			ShowMasterChan:     make(chan *dto.GameShowMasterRequest),
			BottomCardsChan:    make(chan []*dto.Poker),
			PlayCardsChan: []chan []*dto.Poker{
				make(chan []*dto.Poker),
				make(chan []*dto.Poker),
				make(chan []*dto.Poker),
				make(chan []*dto.Poker),
			},
		}
		m.sendGameToClients(game)
		gameResult := game.Run()
		if m.isFinish(game, gameResult) {
			log.Printf("比赛结束!")
			break
		}
		log.Println(gameResult)
		isFristRound = false
		gameBanker = game.Banker
		if gameResult.FinalWinTeam == dto.FirstTeam {
			if gameResult.UpLevel > 0 {
				firstTeamRound = common.NextRoundMap[firstTeamRound]
			}
			if gameBanker == third || gameBanker == fourth {
				gameBanker = banker(turnNextMap[int(gameBanker)])
			}
		} else {
			if gameResult.UpLevel > 0 {
				secondTeamRound = common.NextRoundMap[secondTeamRound]
			}
			if gameBanker == first || gameBanker == second {
				gameBanker = banker(turnNextMap[int(gameBanker)])
			}
		}
		m.sendGameResult(firstTeamRound, secondTeamRound, gameBanker)
	}
}

func (m *Match) sendGameResult(firstTeamRound, secondTeamRound string, banker banker) {
	for i := 0; i < 4; i++ {
		gameResult := dto.GameResultResponse{
			OurRound:      firstTeamRound,
			OtherRound:    secondTeamRound,
			BankerPostion: int(common.GetRelativePos(i, int(banker)-1)),
		}
		if i >= 2 {
			gameResult.OurRound = secondTeamRound
			gameResult.OtherRound = firstTeamRound
		}
		content, _ := json.Marshal(gameResult)
		resp := ResponseMessage{
			MessageType: gameResultResponse,
			Content:     string(content),
		}
		sendMessage, _ := json.Marshal(resp)
		m.PlayerConns[i].WriteMessage(websocket.TextMessage, sendMessage)
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

func (m *Match) isFinish(game *Game, gameResult *dto.GameResult) bool {
	return game.Round == "A" && gameResult.FinalCardWinTeam == game.bankerTeam() && gameResult.Score < 80
}

func (m *Match) sendGameToClients(game *Game) {
	for _, receiveChan := range m.ReceiveGameChans {
		receiveChan <- game
	}
}
