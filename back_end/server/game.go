package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

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
	ShowMasterChan     chan *GameShowMasterRequest
	BottomCards        []*Poker
	BottomCardsChan    chan []*Poker
	PlayCardsChan      []chan []*Poker
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

func (g *Game) sendDealPoker(playerIdx int, poker Poker) {
	content, _ := json.Marshal(poker)
	roomListMessage := ResponseMessage{
		MessageType: dealPoker,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) sendHoleCards(playerIdx int, cards []Poker) {
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

func (g *Game) sendShowMasterPosition(playerIdx int, res *ShowMasterResponse) {
	content, _ := json.Marshal(res)
	roomListMessage := ResponseMessage{
		MessageType: showMasterResult,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(roomListMessage)
	g.PlayerConns[playerIdx].WriteMessage(websocket.TextMessage, sendMessage)
}

func (g *Game) Run() *GameResult {
	defer close(g.ShowMasterChan)
	defer close(g.ShowMasterDoneChan)

	dealPokers := shufflePokers(cardsList)
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
				g.sendShowMasterPosition(i, &ShowMasterResponse{
					Color:              gameReq.Req.Color,
					IsProtect:          gameReq.Req.IsOppose || gameReq.Req.IsSelfProtect,
					IsSelfShowMaster:   i == int(gameReq.Position)-1,
					ShowMasterPosition: int32(getRelativePos(i, int(gameReq.Position)-1)),
				})
			}
		}
	}()
	for i := 0; i < 4; i++ {
		<-g.ShowMasterDoneChan
	}
	log.Println("????????????")
	g.sendHoleCards(int(g.Banker)-1, dealPokers[100:108])
	g.BottomCards = <-g.BottomCardsChan
	for i := 0; i < 8; i++ {
		log.Println(g.BottomCards[i])
	}
	turnPosition := int(g.Banker - 1)
	for {
		for i := 0; i < 4; i++ {
			g.sendPlayTrun(turnPosition)
			cards := <-g.PlayCardsChan[turnPosition]
			for j := 0; j < len(cards); j++ {
				log.Printf("%+v\n", cards[j])
			}
			turnPosition = turnNextMap[turnPosition+1] - 1
		}
	}
	stuckChan := make(chan bool)
	stuckChan <- false
	return &GameResult{}
}

func (g *Game) receiveShowMasterDone() {
	g.ShowMasterDoneChan <- true
}

func (g *Game) receiveBottomCards(cards []*Poker) {
	g.BottomCardsChan <- cards
}

func (g *Game) receivePlayCards(idx int, cards []*Poker) {
	g.PlayCardsChan[idx] <- cards
}

func (g *Game) receiveShowMaster(req *GameShowMasterRequest) {
	g.ShowMasterChan <- req
}

func shufflePokers(src []Poker) []Poker {
	dst := []Poker{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(src)) {
		dst = append(dst, src[i])
	}
	return dst
}
