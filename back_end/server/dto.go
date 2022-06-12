package server

import "encoding/json"

type PlayerMeta struct {
	UUID       string `json:"UUID"`
	PlayerName string `json:"playerName"`
}

func getPlayerMeta(msg string) (PlayerMeta, error) {
	meta := PlayerMeta{}
	err := json.Unmarshal([]byte(msg), &meta)
	return meta, err
}

type ShowMasterRequest struct {
	Master        Poker `json:"master"`
	IsSelfProtect bool  `json:"isSelfProtect"`
	IsOppose      bool  `json:"isOppose"`
}
