package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

const maxRoomPlayer = 4

type Player struct {
	Name string
	Conn *websocket.Conn
}

type NewJoinMsg struct {
	Position int    `json:"position"`
	Name     string `json:"name"`
}

func (p *Player) notifyNewPlayerJoin(position int, name string) {
	content, _ := json.Marshal(NewJoinMsg{
		Position: position,
		Name:     name,
	})
	message := ResponseMessage{
		MessageType: newPlayerJoinRoom,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(message)
	p.Conn.WriteMessage(websocket.TextMessage, sendMessage)
}

func (p *Player) notifyExistPlayers(names []string) {
	if len(names) == 0 {
		return
	}
	content, _ := json.Marshal(names)
	message := ResponseMessage{
		MessageType: existPlayers,
		Content:     string(content),
	}
	sendMessage, _ := json.Marshal(message)
	p.Conn.WriteMessage(websocket.TextMessage, sendMessage)
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
	var existPlayers []string
	for _, player := range r.Players {
		player.notifyNewPlayerJoin(len(r.Players)+1, joinRoomRequest.PlayerName)
		existPlayers = append(existPlayers, player.Name)
	}
	newPlayer := &Player{
		joinRoomRequest.PlayerName,
		joinRoomRequest.Conn,
	}
	newPlayer.notifyExistPlayers(existPlayers)
	r.Players = append(r.Players, newPlayer)
	return true, len(r.Players), r.ClientHandlerChan, ""
}
