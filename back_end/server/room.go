package server

import "github.com/gorilla/websocket"

type Player struct {
	Name string
	Conn *websocket.Conn
}

type Room struct {
	ID      string
	Players []*Player

	ClientHandlerChan chan []byte
}

func (r *Room) RegisterPlayer(joinRoomRequest *JoinRoomRequest) (bool, int, chan []byte) {
	if len(r.Players) == 4 {
		return false, -1, nil
	}
	r.Players = append(r.Players, &Player{
		joinRoomRequest.PlayerName,
		joinRoomRequest.Conn,
	})
	return true, len(r.Players), r.ClientHandlerChan
}
