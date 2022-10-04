package server

import (
	"testing"

	"github.com/YangzhenZhao/shengji/back_end/server/dto"
	"github.com/stretchr/testify/assert"
)

func Test_isTractorsCards(t *testing.T) {
	type args struct {
		firstCards    dto.Cards
		game          *Game
		firstPosition int
	}
	tests := []struct {
		name           string
		args           args
		wantIsTractors bool
	}{
		{
			name: "happy path",
			args: args{
				firstPosition: 3,
				game:          &Game{Round: "8"},
				firstCards: dto.Cards{
					{Color: "heart", Number: "6"},
					{Color: "heart", Number: "6"},
					{Color: "heart", Number: "5"},
					{Color: "heart", Number: "5"},
				},
			},
			wantIsTractors: true,
		},
		{
			name: "happy path - skip round",
			args: args{
				firstPosition: 3,
				game:          &Game{Round: "6"},
				firstCards: dto.Cards{
					{Color: "heart", Number: "7"},
					{Color: "heart", Number: "5"},
					{Color: "heart", Number: "7"},
					{Color: "heart", Number: "5"},
				},
			},
			wantIsTractors: true,
		},
		{
			name: "unhappy path",
			args: args{
				firstPosition: 3,
				game:          &Game{Round: "2"},
				firstCards: dto.Cards{
					{Color: "club", Number: "A"},
					{Color: "club", Number: "9"},
					{Color: "club", Number: "8"},
					{Color: "club", Number: "6"},
				},
			},
			wantIsTractors: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparator := buildComparator(tt.args.game, tt.args.firstCards, tt.args.firstPosition)
			assert.Equal(t, comparator.firstCardsType == tractorsCards, tt.wantIsTractors)
		})
	}
}

func Test_comparator(t *testing.T) {
	type args struct {
		firstCards    dto.Cards
		game          *Game
		firstPosition int
	}
	tests := []struct {
		name            string
		args            args
		otherCards      []dto.Cards
		otherPosition   []int
		wantWinPostions []int
	}{
		{
			name: "happy path - single",
			args: args{
				firstPosition: 3,
				game: &Game{
					Round:           "2",
					FirstTeamRound:  "2",
					SecondTeamRound: "2",
					Banker:          second,
					MasterColor:     "dianmond",
				},
				firstCards: dto.Cards{
					{Color: "dianmond", Number: "3"},
				},
			},
			otherPosition: []int{1, 2, 0},
			otherCards: []dto.Cards{
				{
					{Color: "club", Number: "2"},
				},
				{
					{Color: "dianmond", Number: "10"},
				},
				{
					{Color: "red", Number: "joker"},
				},
			},
			wantWinPostions: []int{1, 1, 0},
		},
		{
			name: "happy path - tractors",
			args: args{
				firstPosition: 3,
				game: &Game{
					Round:           "2",
					FirstTeamRound:  "2",
					SecondTeamRound: "2",
					Banker:          first,
					MasterColor:     "spade",
				},
				firstCards: dto.Cards{
					{Color: "club", Number: "3"},
					{Color: "club", Number: "3"},
					{Color: "club", Number: "4"},
					{Color: "club", Number: "4"},
				},
			},
			otherPosition: []int{1, 2, 0},
			otherCards: []dto.Cards{
				{
					{Color: "club", Number: "K"},
					{Color: "club", Number: "8"},
					{Color: "club", Number: "10"},
					{Color: "club", Number: "6"},
				},
				{
					{Color: "club", Number: "A"},
					{Color: "club", Number: "9"},
					{Color: "club", Number: "8"},
					{Color: "club", Number: "6"},
				},
				{
					{Color: "heart", Number: "7"},
					{Color: "heart", Number: "8"},
					{Color: "heart", Number: "9"},
					{Color: "heart", Number: "J"},
				},
			},
			wantWinPostions: []int{3, 3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparator := buildComparator(tt.args.game, tt.args.firstCards, tt.args.firstPosition)
			for i := 0; i < 3; i++ {
				comparator.addCards(tt.otherCards[i], tt.otherPosition[i])
				assert.Equal(t, tt.wantWinPostions[i], comparator.winPosition)
			}
		})
	}
}
