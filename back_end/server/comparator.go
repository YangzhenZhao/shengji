package server

import (
	"log"
	"sort"

	"github.com/YangzhenZhao/shengji/back_end/server/common"
	"github.com/YangzhenZhao/shengji/back_end/server/dto"
)

type cardsType string

const (
	singleCard    cardsType = "single"
	pairCards     cardsType = "pair"
	tractorsCards cardsType = "tractors"
	flipCards     cardsType = "flip"
)

type Comparator struct {
	game           Game
	winCards       dto.Cards
	isMaster       bool
	winPosition    int
	firstCardsType cardsType
}

func buildComparator(
	game *Game, firstCards dto.Cards,
	firstPosition int,
) *Comparator {
	comparator := &Comparator{
		game:        *game,
		winCards:    firstCards,
		winPosition: firstPosition,
	}
	sort.Sort(firstCards)
	comparator.firstCardsType = comparator.getCardsType(firstCards)
	if comparator.getCardColor(firstCards[0]) == "master" {
		comparator.isMaster = true
	}
	return comparator
}

func (c *Comparator) addCards(cards dto.Cards, position int) {
	if c.isDifferentColors(cards) {
		return
	}
	if c.isMaster && c.getCardColor(cards[0]) != "master" {
		return
	}
	sort.Sort(cards)
	if c.isNewCardsGreater(cards) {
		c.winCards = cards
		c.winPosition = position
	}
}

func (c *Comparator) isNewCardsGreater(cards dto.Cards) bool {
	newCardsType := c.getCardsType(cards)
	switch c.firstCardsType {
	case singleCard:
		return c.isNewSingleCardGreater(cards[0])
	case pairCards:
		return newCardsType == pairCards && c.isNewSingleCardGreater(cards[0])
	case tractorsCards:
		return newCardsType == tractorsCards && c.isNewSingleCardGreater(cards[0])
	case flipCards:
		log.Println("not support filpCards now!!!")
	}
	return false
}

func (c *Comparator) isNewSingleCardGreater(newCard *dto.Poker) bool {
	if newCard == c.winCards[0] {
		return false
	}
	if c.isMaster {
		oldCardMasterLevel := dto.GetCardMasterLevel(c.winCards[0], c.game.Round)
		newCardMasterLevel := dto.GetCardMasterLevel(newCard, c.game.Round)
		if oldCardMasterLevel != newCardMasterLevel {
			return newCardMasterLevel > oldCardMasterLevel
		}
		if newCardMasterLevel == dto.JokerLevelMaster {
			return newCard.Color == "red"
		}
		if newCardMasterLevel == dto.PlayNumberLevelMaster {
			return newCard.Color == c.game.MasterColor
		}
	}
	return dto.CardValueMap[newCard.Number] > dto.CardValueMap[c.winCards[0].Number]
}

func (c *Comparator) getCardsType(cards dto.Cards) cardsType {
	if len(cards) == 1 {
		return singleCard
	}
	if len(cards) == 2 && cards[0] == cards[1] {
		return pairCards
	}
	if c.isTractorsCards(cards) {
		return tractorsCards
	}
	return flipCards
}

func (c *Comparator) isTractorsCards(cards dto.Cards) bool {
	if len(cards) == 0 || len(cards)%2 == 1 || len(cards) < 4 {
		return false
	}
	for i := 2; i < len(cards); i += 2 {
		wantRound := common.NextRoundMap[cards[i-1].Number]
		if wantRound == c.game.Round {
			wantRound = common.NextRoundMap[wantRound]
		}
		if wantRound != cards[i].Number {
			return false
		}
	}
	return true
}

func (c *Comparator) isDifferentColors(cards dto.Cards) bool {
	if len(cards) == 1 {
		return false
	}
	color := c.getCardColor(cards[0])
	for i := 1; i < len(cards); i++ {
		if c.getCardColor(cards[i]) != color {
			return true
		}
	}
	return false
}

func (c *Comparator) getCardColor(card *dto.Poker) string {
	if card.Number == c.game.PlayNumber() ||
		card.Color == "joker" ||
		card.Color == c.game.MasterColor {
		return "master"
	}
	return card.Color
}
