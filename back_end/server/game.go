package server

import (
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

func (g *Game) Run() {
	g.dealCards()
	g.GameResultChan <- &GameResult{}
}

func (g *Game) dealCards() {

}

func shufflePokers(src []Poker) []Poker {
	dst := []Poker{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(src)) {
		dst = append(dst, src[i])
	}
	return dst
}
