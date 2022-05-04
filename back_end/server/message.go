package server

type requestMessageType string
type responseMessageType string

const (
	joinRoom      requestMessageType = "join_room"
	setPlayerName requestMessageType = "set_player_name"
	prepare       requestMessageType = "prepare"

	roomList         responseMessageType = "room_list"
	joinRoomFail     responseMessageType = "join_room_fail"
	existPlayers     responseMessageType = "exists_players"
	hasPlayerPrepare responseMessageType = "has_player_prepare"
)

type RequestMessage struct {
	MessageType requestMessageType `json:"messageType"`
	Content     string             `json:"content"`
}

type ResponseMessage struct {
	MessageType responseMessageType `json:"messageType"`
	Content     string              `json:"content"`
}
