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
