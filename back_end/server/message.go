package server

type requestMessageType string
type responseMessageType string

const (
	joinRoom       requestMessageType = "join_room"
	setPlayerName  requestMessageType = "set_player_name"
	prepare        requestMessageType = "prepare"
	showMasterDone requestMessageType = "show_master_done"
	showMaster     requestMessageType = "show_master"
	kouCards       requestMessageType = "kou_cards"
	playCards      requestMessageType = "play_cards"

	roomList           responseMessageType = "room_list"
	joinRoomFail       responseMessageType = "join_room_fail"
	existPlayers       responseMessageType = "exists_players"
	hasPlayerPrepare   responseMessageType = "has_player_prepare"
	matchBegin         responseMessageType = "match_begin"
	dealPoker          responseMessageType = "deal_poker"
	showMasterResult   responseMessageType = "show_master_result"
	dealHoleCards      responseMessageType = "deal_hole_cards"
	playTurn           responseMessageType = "play_trun"
	showPlayCards      responseMessageType = "show_play_cards"
	increaseScores     responseMessageType = "increase_scores"
	roundEnd           responseMessageType = "round_end"
	biggestPostion     responseMessageType = "biggest_position"
	gameResultResponse responseMessageType = "game_result"
)

type RequestMessage struct {
	MessageType requestMessageType `json:"messageType"`
	Content     string             `json:"content"`
}

type ResponseMessage struct {
	MessageType responseMessageType `json:"messageType"`
	Content     string              `json:"content"`
}
