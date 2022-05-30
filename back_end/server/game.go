package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
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
	Round           string
	FirstTeamRound  string
	SecondTeamRound string
	IsFristRound    bool
	Banker          banker
	PlayerConns     []*websocket.Conn
	GameResultChan  chan *GameResult
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

func (g *Game) Run() {
	var wg sync.WaitGroup
	dealPokers := shufflePokers(cardsList)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		tmpIdx := i
		go func() {
			defer wg.Done()
			for pokerIdx := tmpIdx * 25; pokerIdx < (tmpIdx+1)*25; pokerIdx++ {
				g.sendDealPoker(tmpIdx, dealPokers[pokerIdx])
			}
		}()
	}
	wg.Add(1)
	go func() {
		// TODO: 添加 "亮主" 相关逻辑
		fmt.Println("等待亮主!")
	}()
	wg.Wait()
	g.GameResultChan <- &GameResult{}
}

func shufflePokers(src []Poker) []Poker {
	dst := []Poker{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(src)) {
		dst = append(dst, src[i])
	}
	return dst
}
