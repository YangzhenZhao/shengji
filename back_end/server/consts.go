package server

import (
	"log"
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

type Poker struct {
	Color  string `json:"color"`
	Number string `json:"number"`
}

var numberList []string
var colorList []string
var cardsList []Poker

func init() {
	numberList = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	colorList = []string{spade, heart, club, dianmond}
	for numIdx := 0; numIdx < 13; numIdx++ {
		for colorIdx := 0; colorIdx < 4; colorIdx++ {
			for i := 0; i < 2; i++ {
				cardsList = append(cardsList, Poker{
					Color:  colorList[colorIdx],
					Number: numberList[numIdx],
				})
			}
		}
	}
	for i := 0; i < 2; i++ {
		cardsList = append(cardsList, Poker{
			Color:  black,
			Number: "joker",
		})
		cardsList = append(cardsList, Poker{
			Color:  red,
			Number: "joker",
		})
	}

	log.Println(cardsList)
	randomList := shufflePokers(cardsList)
	log.Println(randomList)
}
