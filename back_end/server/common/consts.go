package common

import (
	"github.com/YangzhenZhao/shengji/back_end/server/dto"
)

const (
	// 黑桃
	spade = "spade"
	// 红桃
	heart = "heart"
	// 梅花
	club = "club"
	// 方块
	dianmond = "dianmond"
	red      = "red"
	black    = "black"
)

var TeamMateMap = map[int]int{
	0: 1,
	1: 0,
	2: 3,
	3: 2,
}
var OpponentMap = map[int][2]int{
	0: {2, 3},
	1: {3, 2},
	2: {1, 0},
	3: {0, 1},
}

var numberList []string
var colorList []string
var CardsList []dto.Poker

func init() {
	numberList = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	colorList = []string{spade, heart, club, dianmond}
	for numIdx := 0; numIdx < 13; numIdx++ {
		for colorIdx := 0; colorIdx < 4; colorIdx++ {
			for i := 0; i < 2; i++ {
				CardsList = append(CardsList, dto.Poker{
					Color:  colorList[colorIdx],
					Number: numberList[numIdx],
				})
			}
		}
	}
	for i := 0; i < 2; i++ {
		CardsList = append(CardsList, dto.Poker{
			Color:  black,
			Number: "joker",
		})
		CardsList = append(CardsList, dto.Poker{
			Color:  red,
			Number: "joker",
		})
	}
}
