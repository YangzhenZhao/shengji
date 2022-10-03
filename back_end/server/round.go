package server

import (
	"log"

	"github.com/YangzhenZhao/shengji/back_end/server/common"
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

// TODO: 支持甩牌
// 目前只支持单牌，对子，拖拉机
func (r *Round) run() *dto.RoundResult {
	turnPosition := r.firstTurnPosition
	cardsTotalScore := 0
	var comparator *Comparator
	var cardsLength int
	for i := 0; i < 4; i++ {
		r.game.sendPlayTrun(turnPosition)
		cards := <-r.game.PlayCardsChan[turnPosition]
		cardsLength = len(cards)
		cardsTotalScore += getCardsScores(cards)
		for j := 0; j < len(cards); j++ {
			log.Printf("%+v\n", cards[j])
		}
		for j := 0; j < 4; j++ {
			if j != turnPosition {
				r.game.sendShowPlayCards(turnPosition, j, cards)
			}
		}
		if i == 0 {
			comparator = buildComparator(r.game, cards, turnPosition)
		} else {
			comparator.addCards(cards, turnPosition)
		}
		for j := 0; j < 4; j++ {
			r.game.sendBiggestPostion(comparator.winPosition, j)
		}
		turnPosition = turnNextMap[turnPosition+1] - 1
	}
	increaseScore := 0
	if r.isEarnScoreTeamWin(comparator) {
		increaseScore = cardsTotalScore
	}
	return &dto.RoundResult{
		WinPosition:   comparator.winPosition,
		IncreaseScore: increaseScore,
		CardsLength:   cardsLength,
	}
}

func (r *Round) isEarnScoreTeamWin(comparator *Comparator) bool {
	return comparator.winPosition+1 != int(r.game.Banker) &&
		common.TeamMateMap[comparator.winPosition]+1 != int(r.game.Banker)
}

func getCardsScores(cards []*dto.Poker) int {
	sumScore := 0
	for _, card := range cards {
		if score, ok := common.ScoreMap[card.Number]; ok {
			sumScore += score
		}
	}
	return sumScore
}
