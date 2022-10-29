package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/YangzhenZhao/shengji/back_end/server/common"
	"github.com/YangzhenZhao/shengji/back_end/server/dto"
	"github.com/gorilla/websocket"
)

type banker int

const (
	unknown banker = iota
	first
	second
	third
	fourth
)

type Game struct {
	Round              string
	FirstTeamRound     string
	SecondTeamRound    string
	IsFristRound       bool
	Banker             banker
	MasterColor        string
	IsProtect          bool
	PlayerConns        []*websocket.Conn
	ShowMasterDoneChan chan bool
	ShowMasterChan     chan *dto.GameShowMasterRequest
	BottomCards        []*dto.Poker
	BottomCardsChan    chan []*dto.Poker
	PlayCardsChan      []chan []*dto.Poker
}

var turnNextMap = map[int]int{
	1: 4,
	4: 2,
	2: 3,
	3: 1,
}

func (g *Game) bankerTeam() dto.TeamType {
	if g.Banker == unknown {
		return dto.UnknownTeam
	}
	if g.Banker == first || g.Banker == second {
		return dto.FirstTeam
	}
	return dto.SecondTeam
}

func (g *Game) sendDealPoker(playerIdx int, poker dto.Poker) {
	content, _ := json.Marshal(poker)
	roomListMessage := ResponseMessage{
		MessageType: dealPoker,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) sendHoleCards(playerIdx int, cards []dto.Poker) {
	content, _ := json.Marshal(cards)
	roomListMessage := ResponseMessage{
		MessageType: dealHoleCards,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) sendPlayTrun(playerIdx int) {
	roomListMessage := ResponseMessage{
		MessageType: playTurn,
		Content:     "",
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) sendShowPlayCards(showCardsIdx, playerIdx int, cards []*dto.Poker) {
	relativeIdx := common.GetRelativePos(playerIdx, showCardsIdx)
	content, _ := json.Marshal(dto.ShowPlayCardsResponse{
		ShowIdx: relativeIdx,
		Cards:   cards,
	})
	resp := ResponseMessage{
		MessageType: showPlayCards,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(resp)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) sendBiggestPostion(biggestPostionIdx, playerIdx int) {
	relativeIdx := common.GetRelativePos(playerIdx, biggestPostionIdx)
	resp := ResponseMessage{
		MessageType: biggestPostion,
		Content:     strconv.Itoa(relativeIdx),
	}
	sendMessage, _ := json.Marshal(resp)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) sendIncreaseScores(score int) {
	resp := ResponseMessage{
		MessageType: increaseScores,
		Content:     strconv.Itoa(score),
	}
	sendMessage, _ := json.Marshal(resp)
	for i := 0; i < 4; i++ {
		g.PlayerConns[i].WriteMessage(websocket.TextMessage, sendMessage)
	}
}

func (g *Game) sendRoundEnd() {
	resp := ResponseMessage{
		MessageType: roundEnd,
	}
	sendMessage, _ := json.Marshal(resp)
	for i := 0; i < 4; i++ {
		g.PlayerConns[i].WriteMessage(websocket.TextMessage, sendMessage)
	}
}

func (g *Game) sendShowMasterPosition(playerIdx int, res *dto.ShowMasterResponse) {
	content, _ := json.Marshal(res)
	roomListMessage := ResponseMessage{
		MessageType: showMasterResult,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) Run() *dto.GameResult {
	defer close(g.ShowMasterChan)
	defer close(g.ShowMasterDoneChan)

	dealPokers := shufflePokers(common.CardsList)
	for i := 0; i < 4; i++ {
		tmpIdx := i
		go func() {
			for pokerIdx := tmpIdx * 25; pokerIdx < (tmpIdx+1)*25; pokerIdx++ {
				g.sendDealPoker(tmpIdx, dealPokers[pokerIdx])
			}
		}()
	}
	go func() {
		for gameReq := range g.ShowMasterChan {
			log.Println("receive game req.........")
			if g.IsFristRound {
				g.Banker = banker(gameReq.Position)
			}
			g.MasterColor = gameReq.Req.Color
			if gameReq.Req.IsSelfProtect || gameReq.Req.IsOppose {
				g.IsProtect = true
			}
			for i := 0; i < 4; i++ {
				g.sendShowMasterPosition(i, &dto.ShowMasterResponse{
					Color:              gameReq.Req.Color,
					IsProtect:          gameReq.Req.IsOppose || gameReq.Req.IsSelfProtect,
					IsSelfShowMaster:   i == int(gameReq.Position)-1,
					ShowMasterPosition: int32(common.GetRelativePos(i, int(gameReq.Position)-1)),
				})
			}
		}
	}()
	for i := 0; i < 4; i++ {
		<-g.ShowMasterDoneChan
	}
	log.Println("亮主完成")
	g.sendHoleCards(int(g.Banker)-1, dealPokers[100:108])
	g.BottomCards = <-g.BottomCardsChan
	for i := 0; i < 8; i++ {
		log.Println(g.BottomCards[i])
	}
	turnPosition := int(g.Banker - 1)
	scores := 0
	playTotalCardsNum := 0
	finalCardWinTeam := dto.FirstTeam
	for {
		round := newRound(g, turnPosition)
		roundResult := round.run()
		turnPosition = roundResult.WinPosition
		scores += roundResult.IncreaseScore
		if roundResult.IncreaseScore > 0 {
			g.sendIncreaseScores(roundResult.IncreaseScore)
		}
		time.Sleep(2 * time.Second)
		g.sendRoundEnd()
		playTotalCardsNum += roundResult.CardsLength
		if playTotalCardsNum == common.PlayerInitCardsNumber {
			if roundResult.WinPosition > 1 {
				finalCardWinTeam = dto.SecondTeam
			}
			break
		}
	}
	upLevel := 1
	var finalWinTeam dto.TeamType
	if g.Banker == first || g.Banker == second {
		finalWinTeam = dto.FirstTeam
		if finalCardWinTeam == dto.SecondTeam || scores >= 80 {
			if finalCardWinTeam == dto.FirstTeam || scores < 80 {
				upLevel = 0
			}
			finalWinTeam = dto.SecondTeam
		}
	} else {
		finalWinTeam = dto.SecondTeam
		if finalCardWinTeam == dto.FirstTeam || scores >= 80 {
			if finalCardWinTeam == dto.SecondTeam || scores < 80 {
				upLevel = 0
			}
			finalWinTeam = dto.FirstTeam
		}
	}

	return &dto.GameResult{
		Score:            scores,
		FinalCardWinTeam: finalCardWinTeam,
		FinalWinTeam:     finalWinTeam,
		UpLevel:          upLevel,
	}
}

func (g *Game) receiveShowMasterDone() {
	g.ShowMasterDoneChan <- true
}

func (g *Game) PlayNumber() string {
	if g.Banker == first || g.Banker == second {
		return g.FirstTeamRound
	}
	return g.SecondTeamRound
}

func (g *Game) receiveBottomCards(cards []*dto.Poker) {
	g.BottomCardsChan <- cards
}

func (g *Game) receivePlayCards(idx int, cards []*dto.Poker) {
	g.PlayCardsChan[idx] <- cards
}

func (g *Game) receiveShowMaster(req *dto.GameShowMasterRequest) {
	g.ShowMasterChan <- req
}

func shufflePokers(src []dto.Poker) []dto.Poker {
	dst := []dto.Poker{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(src)) {
		dst = append(dst, src[i])
	}
	return dst
}
