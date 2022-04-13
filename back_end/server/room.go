package server

import "github.com/gorilla/websocket"

const maxRoomPlayer = 4

type Player struct {
	Name string
	Conn *websocket.Conn
}

type Room struct {
	ID      string
	Players []*Player

	ClientHandlerChan chan []byte
}

func (r *Room) RegisterPlayer(joinRoomRequest *JoinRoomRequest) (bool, int, chan []byte, string) {
	for _, player := range r.Players {
		if player.Name == joinRoomRequest.PlayerName {
			return false, -1, nil, "不可重复加入!"
		}
	}
	if len(r.Players) == maxRoomPlayer {
		return false, -1, nil, "房间已满!"
	}
	r.Players = append(r.Players, &Player{
		joinRoomRequest.PlayerName,
		joinRoomRequest.Conn,
	})
	return true, len(r.Players), r.ClientHandlerChan, ""
}
