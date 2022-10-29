package dto

import (
	"encoding/json"
)

type Poker struct {
	Color  string `json:"color"`
	Number string `json:"number"`
}

var CardValueMap = map[string]int{
	"2":  2,
	"3":  3,
	"4":  4,
	"5":  5,
	"6":  6,
	"7":  7,
	"8":  8,
	"9":  9,
	"10": 10,
	"J":  11,
	"Q":  12,
	"K":  13,
	"A":  14,
}

type Cards []*Poker

func (c Cards) Len() int {
	return len(c)
}

func (c Cards) Less(i, j int) bool {
	return CardValueMap[c[i].Number] < CardValueMap[c[j].Number]
}

func (c Cards) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type PlayerMeta struct {
	UUID       string `json:"UUID"`
	PlayerName string `json:"playerName"`
}

// GetPlayerMeta ...
func GetPlayerMeta(msg string) (PlayerMeta, error) {
	meta := PlayerMeta{}
	err := json.Unmarshal([]byte(msg), &meta)
	return meta, err
}

type ShowMasterRequest struct {
	Color         string `json:"color"`
	IsSelfProtect bool   `json:"isSelfProtect"`
	IsOppose      bool   `json:"isOppose"`
}

type GameShowMasterRequest struct {
	Req      ShowMasterRequest
	Position int32
}

type ShowMasterResponse struct {
	Color              string `json:"color"`
	IsProtect          bool   `json:"isProtect"`
	IsSelfShowMaster   bool   `json:"isSelfShowMaster"`
	ShowMasterPosition int32  `json:"showMasterPosition"`
}

type ShowPlayCardsResponse struct {
	ShowIdx int      `json:"showIdx"`
	Cards   []*Poker `json:"cards"`
}

type GameResultResponse struct {
	OurRound      string `json:"ourRound"`
	OtherRound    string `json:"otherRound"`
	BankerPostion int    `json:"bankerPostion"`
}
