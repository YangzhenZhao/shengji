package server

import (
	"log"

	"github.com/YangzhenZhao/shengji/back_end/server/dto"
)

type Round struct {
	game              *Game
	firstTurnPosition int
}

func newRound(game *Game, firstTurnPosition int) *Round {
	return &Round{
		game:              game,
		firstTurnPosition: firstTurnPosition,
	}
}

func (r *Round) run() *dto.RoundResult {
	turnPosition := r.firstTurnPosition
	for i := 0; i < 4; i++ {
		r.game.sendPlayTrun(turnPosition)
		cards := <-r.game.PlayCardsChan[turnPosition]
		for j := 0; j < len(cards); j++ {
			log.Printf("%+v\n", cards[j])
		}
		for j := 0; j < 4; j++ {
			if j != turnPosition {
				r.game.sendShowPlayCards(turnPosition, j, cards)
			}
		}
		turnPosition = turnNextMap[turnPosition+1] - 1
	}
	return &dto.RoundResult{
		WinPosition: 0,
	}
}
