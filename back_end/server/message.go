package server

type requestMessageType string
type responseMessageType string

const (
	joinRoom      requestMessageType = "join_room"
	setPlayerName requestMessageType = "set_player_name"

	roomList     responseMessageType = "room_list"
	joinRoomFail responseMessageType = "join_room_fail"
)

type RequestMessage struct {
	MessageType requestMessageType `json:"messageType"`
	Content     string             `json:"content"`
}

type ResponseMessage struct {
	MessageType responseMessageType `json:"messageType"`
	Content     string              `json:"content"`
}
