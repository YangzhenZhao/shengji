package server

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RegisterClientRequest struct {
	PlayerName string
	Client     *Client
}

type JoinRoomRequest struct {
	PlayerName string
	RoomID     string
	Conn       *websocket.Conn
}

type Hub struct {
	Clients map[string]*Client
	Rooms   map[string]*Room

	RegisterClientChan  chan *RegisterClientRequest
	JoinRoomRequestChan chan *JoinRoomRequest
}

func NewHub() *Hub {
	roomID := uuid.NewString()
	return &Hub{
		Clients: make(map[string]*Client),
		Rooms: map[string]*Room{
			roomID: {
				ID:      roomID,
				Players: []*Player{},
			},
		},
		RegisterClientChan:  make(chan *RegisterClientRequest),
		JoinRoomRequestChan: make(chan *JoinRoomRequest),
	}
}

func (h *Hub) roomList() []string {
	var res []string
	for roomID := range h.Rooms {
		res = append(res, roomID)
	}
	return res
}

func (h *Hub) Run() {
	for {
		select {
		case registerClientRequest := <-h.RegisterClientChan:
			h.Clients[registerClientRequest.PlayerName] = registerClientRequest.Client
			log.Printf("%s 注册成功!\n", registerClientRequest.PlayerName)
		case joinRoomRequest := <-h.JoinRoomRequestChan:
			room := h.Rooms[joinRoomRequest.RoomID]
			isJoin, position, clientHandlerChan := room.RegisterPlayer(joinRoomRequest)
			if isJoin {
				client := h.Clients[joinRoomRequest.PlayerName]
				client.Room = &RoomDetail{
					Position:    position,
					HandlerChan: clientHandlerChan,
				}
				log.Printf("玩家 %s 加入房间成功! 房间 ID: %s\n", joinRoomRequest.PlayerName, joinRoomRequest.RoomID)
			}
		}
	}
}
