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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparator := buildComparator(tt.args.game, tt.args.firstCards, tt.args.firstPosition)
			assert.Equal(t, comparator.firstCardsType == tractorsCards, tt.wantIsTractors)
		})
	}
}
