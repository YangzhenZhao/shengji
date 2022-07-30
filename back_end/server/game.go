package server

import (
	"encoding/json"
	"log"
	"math/rand"
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

type teamType int

const (
	unknownTeam teamType = iota
	firstTeam
	secondTeam
)

type GameResult struct {
	Score            int
	FinalCardWinTeam teamType
}

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

func (g *Game) bankerTeam() teamType {
	if g.Banker == unknown {
		return unknownTeam
	}
	if g.Banker == first || g.Banker == second {
		return firstTeam
	}
	return secondTeam
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

func (g *Game) sendShowMasterPosition(playerIdx int, res *dto.ShowMasterResponse) {
	content, _ := json.Marshal(res)
	roomListMessage := ResponseMessage{
		MessageType: showMasterResult,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) calcWinPos(turnCards [][]*dto.Poker) int {
	return 0
}

func (g *Game) Run() *GameResult {
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
			g.Banker = banker(gameReq.Position)
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
	for {
		turnCards := [][]*dto.Poker{
			{},
			{},
			{},
			{},
		}
		for i := 0; i < 4; i++ {
			g.sendPlayTrun(turnPosition)
			cards := <-g.PlayCardsChan[turnPosition]
			turnCards[turnPosition] = cards
			for j := 0; j < len(cards); j++ {
				log.Printf("%+v\n", cards[j])
			}
			for j := 0; j < 4; j++ {
				if j != turnPosition {
					g.sendShowPlayCards(turnPosition, j, cards)
				}
			}
			turnPosition = turnNextMap[turnPosition+1] - 1
		}
		turnPosition = g.calcWinPos(turnCards)
	}
	stuckChan := make(chan bool)
	stuckChan <- false
	return &GameResult{}
}

func (g *Game) receiveShowMasterDone() {
	g.ShowMasterDoneChan <- true
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
